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
case class Food extends Tile
case object Rock extends Tile
case object Water extends Tile
case object Sugar extends Food
case object Mill extends Food
case object Meat extends Food


/**
 * GameInfo holds and updates game information such as:
 * - ants on the battleground
 * - map
 */
object GameInfo 
{
  var turnId: Int = 0
  var antsPerPlayer: Int = 0
  var nbPlayers: Int = 0
  var playing: Boolean = false
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

  def update = {
    readInfo
    act
  }

  def readInfo = {
    // ID X Y DX DY Brain
    val pattern = "[0-9]+".r
  }

  def act = state match {
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
      ai.act
      ai.changeState(Retreat).act
    }

}
