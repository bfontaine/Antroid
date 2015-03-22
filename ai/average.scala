/**
 * Ant holds all information we can have about an ant
 */ 
      // N lines: X Y (C)ontent (S)ee_now
class Ant
{
  // id
  var id: Int = 0
  // position
  var x: Int = 0
  var y: Int = 0
  // direction
  var dx: Int = 0
  var dy: Int = 0
  var energy: Int = 0
  var acid: Int = 0
  var brain: Int = 0

  def init (_id: Int, _x: Int, _y: Int, _dx: Int, _dy: Int, _e: Int, _a: Int, _b: Int) : Ant = 
    {
      id = _id
      x = _x
      y = _y
      dx = _dx
      dy = _dy
      energy = _e
      acid = _a
      brain = _b
      return this
    }

  def init (_x: Int, _y: Int, _dx: Int, _dy: Int, _b: Int) : Ant = 
    {
      x = _x
      y = _y
      dx = _dx
      dy = _dy
      brain = _b
      return this
    }

}

class Point (val x: Int, val y: Int)

object Point 
{
  def fromCoord (x: Int, y: Int) = new Point(x,y)
}

/**
 * Tile is the enumeration of all possible field value, used in the map
 */
sealed abstract class Tile
case object Grass extends Tile
case object Rock extends Tile
case object Water extends Tile
case object Sugar extends Tile
case object Mill extends Tile
case object Meat extends Tile

object Tile 
{
  class UnknownTile extends Exception
  def fromInt (i: Int) = i match {
    case 0 => Grass;
    case 1 => Sugar;
    case 2 => Rock;
    case 3 => Mill;
    case 4 => Water;
    case 5 => Meat;
    case _ => throw new UnknownTile;
  }
}

class VisibleTile (val tile: Tile, var visible: Boolean = false);

/**
 * GameInfo holds and updates game information such as:
 * - ants on the battleground
 * - map
 */
object GameInfo 
{
  // general game information
  private var _turnId: Int = 0
  def turnId = _turnId
  private var _antsPerPlayer: Int = 0
  def antsPerPlayer = _antsPerPlayer
  private var _nbPlayers: Int = 0
  def nbPlayers = _nbPlayers
  private var _playing: Boolean = false
  def playing = _playing
  // has the game initialised (all objects have already been created)
  private var initialised = false

  // player's ants
  private var _myAnts : Array[Ant] = Array ()
  def myAnts = _myAnts

  // other ants we can see at this turn
  private var _enemyAnts : List[Ant] = List()
  // other ants we could see on previous turn
  private var _oldEnemies : List[Ant] = List()

  // map
  private var _map : Map[Point, VisibleTile] = Map()
  private var _width : Int = 0
  private var _height : Int = 0

  /*
   * Tries to read information on stdin in order to update state
   */ 
  def nextTurn () = 
    {
      // Protocol: 
      // (T)urn_id (A)nts/player (P)layersNb (S)tatus
      // A lines for ants: ID X Y DX DY (E)nergy (A)cid (B)rain
      // (N)b_enemies
      // N lines for opponents' ants: X Y DX DY B
      // Map: W H (N)b_points
      // N lines: X Y (C)ontent (S)ee_now
      var finished = false
      var init = false
      while (! finished) 
      {
	var line = nullexception(readLine())
	// header
	header (line)
	// player's ants treatment
	ants (_antsPerPlayer)
	// opponents' ants header + treatment
	opponentAnts ()
	// map header + treatment
	map ()
      }
    }

  val header_pattern = "([0-9]+) ([0-9]+) ([0-9]+) ([0-1]+)".r
  private def header (line: String) = 
    {
      // (T)urn_id (A)nts/player (P)layersNb (S)tatus
      val header_pattern(t, a, p, s) = line
      _playing = if (s == "0") {false} else {true}
      _turnId = t.toInt
      if (! initialised) 
	init (a.toInt, p.toInt)
    }

  private def init (ants: Int, players: Int) = 
    {
      _antsPerPlayer = ants
      _nbPlayers = players
      _myAnts = Array.fill(ants){new Ant}
    }

  private def ants (n: Int) = 
    {
      for (_ <- 0 until n) {
	val line = nullexception(readLine())
	ant (line)
      }
    }
  
  val ant_pattern = "([0-9]+) ([0-9]+) ([0-9]+) (-?[0-9]+) (-?[0-9]+) ([0-9]+) ([0-9]+) ([0-9]+)".r
  private def ant (line: String) = 
    {
      // ID X Y DX DY (E)nergy (A)cid (B)rain
      val ant_pattern(id, x, y, dx, dy, e, a, b) = line
      val ant = (new Ant).init(id.toInt, x.toInt, y.toInt, dx.toInt, dy.toInt, e.toInt, a.toInt, b.toInt)
      _myAnts(id.toInt) = ant
    }

  private def opponentAnts () = 
    {
      // flushing 
      _oldEnemies = _enemyAnts
      _enemyAnts = List()
      // treatment
      val s = readLine()
      val n = s.toInt
      for (_ <- 0 until n) {
	val line = nullexception(readLine())
	opponentAnt (line)
      }
    }

  val opponent_pattern = "([0-9]+) ([0-9]+) ([0-9]+) ([0-9]+) ([0-9]+)".r
  private def opponentAnt (line: String) = 
    {
      // X Y DX DY B
      val opponent_pattern(x, y, dx, dy, b) = line
      val ant = (new Ant).init(x.toInt, y.toInt, dx.toInt, dy.toInt, b.toInt)
      _enemyAnts = _enemyAnts :+ ant      
    }

  private def map () = 
    {
      // resetting map visibility
      for ((p,t) <- _map) {
	t.visible = false
      }
      // parsing
      val line = nullexception(readLine())      
      val nTiles = mapHeader (line)
      for (i <- 0 until nTiles) 
	{
	  val line = nullexception(readLine())
	  mapTiles (line)
	}
    }

  val mapHeader_pattern = "([0-9]+) ([0-9]+) ([0-9]+)".r
  private def mapHeader (line: String) : Int = 
    {
      // W H (N)b_points
      val mapHeader_pattern(w, h, n) = line
      _width = w.toInt
      _height = h.toInt
      return n.toInt
    }

  private val mapTiles_pattern = "([0-9]+) ([0-9]+) ([0-9]+)".r
  private def mapTiles (line: String) = 
    {
      // X Y (C)ontent (S)ee_now
      val mapTiles_pattern(x, y, c, s) = line
      val point = Point.fromCoord(x.toInt, y.toInt)
      val tile = _map.get(point)
      tile match {
	case Some(t) => t.visible = true; _map = _map + (point -> t);
	case None => val t = new VisibleTile(Tile.fromInt(c.toInt),true);
	_map = _map+(point -> t)
      }
    }

  class BadServerPacket extends Exception

  /**
   * checks if [line] is null (it would mean that reading stdin results in EOF
   * while it shouldn't).
   * It returns the same string to be fluent. 
   */ 
  private def nullexception (line: String) : String = 
    {
      if (line == null)
	throw new BadServerPacket
      else 
	return line
    }

}

/**
 * A State is affected to each AI agent. It will act differently according
 * to the state, which can be changed at any time if required
 */
sealed abstract class AIState
// the Ant should just wait without moving
case object Wait extends AIState
// the Ant should explore unknown tiles 
case object Explore extends AIState 
// the Ant should reach the closest ally 
case object Retreat extends AIState
// the Ant should try to attack closest opponent 
case object Battle extends AIState 
// the Ant should reach the ally given as parameter
case class Unite (lead: Int) extends AIState

class AIAgent
{
  var state : AIState = Wait

  def act () = state match {
    case Wait => do_wait;
    case Explore => do_explore;
    case Retreat => do_retreat;
    case Battle => do_battle;
    case Unite(n) => do_unite(n);
    case _ => println("Unknown state")
  }

  def changeState (s: AIState) = { state = s; this }

  private def do_wait = println("I am waiting")

  private def do_explore = println("I am exploring")

  private def do_retreat = println("I am retreating")

  private def do_battle = println("I want to fight")

  private def do_unite (lead: Int) = println("I should join $lead")

}

object Test {

  def main (args: Array[String]) = 
    {
      val ai = new AIAgent
      ai.act()
      ai.changeState(Retreat).act
    }

}
