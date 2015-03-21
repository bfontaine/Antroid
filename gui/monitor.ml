(* Compile with
 * [ocamlopt -w -8 str.cmxa graphics.cmxa monitor.ml] *)

open Graphics

(* Helper applying [fn ()] [n] times and return the list of results. *)
let loop fn n =
  let rec loop acc fn n =
    if n = 0 then acc
    else loop (fn () :: acc) fn (n - 1) in
  loop [] fn n

(* Helper reading a list of integers separated by space on stdin. *)
let read_ints () = input_line stdin
                   |> Str.split (Str.regexp_string " ")
                   |> List.map int_of_string

let color_of_content = function
  | 0 -> rgb  51 102   0 (* grass: green  *)
  | 2 -> rgb 192 192 192 (*  rock: grey   *)
  | 4 -> rgb   0 153 153 (* water: blue   *)
  | 1 -> rgb 255 255 255 (* sugar: white  *)
  | 3 -> rgb 255 255 102 (* wheat: yellow *)
  | 5 -> rgb 255 153 151 (*  meat: brown  *)
  | _ -> rgb 255   0 127 (*     ?: pink   *)

let _ =

  (* Choose the screen size on the command line or use default. *)
  open_graph (if Array.length Sys.argv = 1 then "" else Sys.argv.(1)) ;

  (* Yes, we are setting the title of the window. *)
  set_window_title "Antroid" ;

  (* Disable the automatic synchronization in order to draw the whole
   * map before refreshing screen. *)
  auto_synchronize false ;

  (* Loop until the game is over. *)
  while true do

    let [t; a; p; s] = read_ints () in

    (* If game is over, wait for user to press a key and exit. *)
    if s = 0 then begin ignore (wait_next_event [Key_pressed]) ;
                        exit 0 end ;

    (* Read ally ants info.  *)
    let ants = loop read_ints a in

    (* Read visible enemy ants info. *)
    let [n] = read_ints () in
    let enemies = loop read_ints n in

    (* Read the map description. *)
    let [w ; h ; n] = read_ints () in
    let map = loop read_ints n in

    (* Get the window size and compute the size of a cell. *)
    let scr_w = size_x () in
    let scr_h = size_y () in
    let side = min (scr_w / w) (scr_h / h) in

    (* Draw a black rectangle, otherwise background is white. *)
    set_color black ; fill_rect 0 0 scr_w scr_h ;

    (* Helper: draw a square at [(x, y)]. *)
    let square x y = fill_rect (x * side) (y * side) side side in

    (* Draw map. *)
    List.iter (fun [x ; y; c; _] -> set_color (color_of_content c) ;
                                    square x y) map ;

    (* Allies *)
    set_color blue ;
    List.iter (fun [_; x; y; _; _; _; _; _] -> square x y) ants ;

    (* Enemies *)
    set_color blue ;
    List.iter (fun [x; y; _; _; _] -> square x y) enemies ;

    (* Refresh the screen. *)
    synchronize ()

  done
