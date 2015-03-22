/**
 * Ant holds all information we can have about an ant
 */ 
class Ant (val id: Int)
{
  // position
  var x: Int = 0
  var y: Int = 0
  // direction
  var dx: Int = 0
  var dy: Int = 0
  var energy: Int = 0
  var acid: Int = 0
  var brain: Int = 0
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


/**
 * GameInfo holds and updates game information such as:
 * - ants on the battleground
 * - map
 */
object GameInfo 
{
  private var _turnId: Int = 0
  def turnId = _turnId
  private var _antsPerPlayer: Int = 0
  def antsPerPlayer = _antsPerPlayer
  private var _nbPlayers: Int = 0
  def nbPlayers = _nbPlayers
  private var _playing: Boolean = false
  def playing = _playing

  /*
   * Tries to read information on stdin in order to update state
   */ 
  def nextTurn () = 
    {
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
      val header_pattern(_turnId, _antsPerPlayer, _nbPlayers, status) = line
      _playing = if (status == "0") {false} else {true}
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
      val ant_pattern(id, x, y, dx, dy, e, a, b) = line
    }

  private def opponentAnts () = 
    {
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
      val opponent_pattern(x, y, dx, dy, b) = line
    }

  private def map () = 
    {
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
      val mapHeader_pattern(w, h, n) = line
      return n.toInt
    }

  private val mapTiles_pattern = "([0-9]+) ([0-9]+) ([0-9]+)".r
  private def mapTiles (line: String) = 
    {
      val mapTiles_pattern(x, y, c, s) = line
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
    case Wait => do_wait
      case Explore => do_explore
	case Retreat => do_retreat
	  case Battle => do_battle
	    case Unite(n) => do_unite(n)
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
