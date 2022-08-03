
let file_opts = {
    types: [{
      description: 'Yaml file',
      accept: {'text/plain': ['.yml', '.yaml']},
    }],
};
let project = {};
let currentRoom = null;
let toggleButton = document.getElementById("toggleEditor");
toggleButton.onclick = function(event){
    event.preventDefault();
    if (map.$el.style.visibility == 'visible') {
        show_editor();
    } else {
        hide_editor();
    }
    
}
function map_room_clicked(room) {
    show_room_in_editor(room.id);
    map.show(project, room);
}
function popup_close() {
    let popup = document.getElementById("popup");
    popup.style.visibility = 'hidden';
    popup.getElementsByClassName("popup-content")[0].innerHTML="";
    return false;
}
function help_about() {
    let popup = document.getElementById("popup");
    popup.style.visibility = 'visible';
    popup.getElementsByClassName("popup-content")[0].innerHTML="About<br/><br/>SWR Editor<br/>&copy; 2022; All rights reserved.";
    return false;
}
function show_editor() {
    let $editor = document.getElementById("editor");
    $editor.style.visibility = 'visible';
    $editor.style.display = 'block';
    map.$el.style.visibility = 'hidden';
    map.$el.style.display = 'none';
    editor.resize();
}
function hide_editor() {
    let $editor = document.getElementById("editor");
    $editor.style.visibility = 'hidden';
    $editor.style.display = 'none';
    map.$el.style.visibility = 'visible';
    map.$el.style.display = 'block';
    map.show(project, currentRoom);
}
function make_tree_node(contents, index) {
    let $node = document.createElement("li");
    $node.innerHTML = contents;
    if (index != undefined) {
        $node.setAttribute("index", index);
    }
    return $node;
}
function make_tree() {
    let $rooms = document.createElement("ul");
    $rooms.style.cursor = "pointer";
    for (const room of project.rooms) {
        $n = make_tree_node(room.id);
        $n.style.cursor = "pointer";
        $n.onclick = function(event){ 
            currentRoom = room;
            event.preventDefault();
            show_room_in_editor(room.id);
            map.show(project, room);
        }
        $rooms.appendChild($n);
    }
    return $rooms;
}
function show_room_in_editor(roomId) {
    console.log(editor.find(`id: ${roomId}`, {range: null}, true));
    editor.focus();
}
function show_project_tree() {
    let $sidePanel = document.getElementsByClassName("sidepanel")[0];
    $sidePanel.innerHTML = "";

    let $root = document.createElement("div");
    $root.id = "treeview";
    $root.appendChild(make_tree());
    $sidePanel.appendChild($root);
}
function new_file() {
    show_editor("");
    return false;
}

async function open_file() {
    [fileHandle] = await window.showOpenFilePicker(file_opts);
    const file = await fileHandle.getFile();
    const contents = await file.text();
    project = YAML.parse(contents);
    show_project_tree();
    editor.session.setValue(contents);
    //show_editor(contents);
    currentRoom = project.rooms[0];
    map.show(project, project.rooms[0]);
    filename = file.name;
    return true;
}

async function save_file() {
    let contents = editor.getValue();
    let fileHandle = await window.showSaveFilePicker(file_opts);
    const writable = await fileHandle.createWritable();
    // Write the contents of the file to the stream.
    await writable.write(contents);

    // Close the file and write the contents to disk.
    await writable.close();
    return true;
}

function quit() {
    window.close();
}

function create_star_system() {

}

function create_area() {

}

function create_entity() {

}

function create_item() {

}

function create_ship() {

}

function create_shop() {

}

function find_player() {

}

function find_star_system() {

}

function find_area() {

}

function find_entity() {

}

function find_item() {

}

function find_ship() {

}

function find_shop() {

}

function server_options() {

}

function server_preferences() {

}

function server_restrict() {

}

function server_restart() {

}

function server_shutdown() {

}

function server_logs() {

}

let editor = ace.edit("editor");
    editor.setOption("useSoftTabs", true);
    editor.setTheme("ace/theme/terminal");
    editor.setShowPrintMargin(false);
    editor.session.setMode("ace/mode/yaml");
    editor.resize();

let map = mapper().new(document.getElementById("map"), map_room_clicked);