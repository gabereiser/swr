
let file_opts = {
    types: [{
      description: 'Yaml file',
      accept: {'text/plain': ['.yml', '.yaml']},
    }],
};
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

function new_file() {
    let $editor = document.getElementById("editor");
    $editor.style.visibility = 'visible';
    $editor.style.display = 'block';
    editor.session.setValue("");
    editor.resize();
    return false;
}

async function open_file() {
    [fileHandle] = await window.showOpenFilePicker(file_opts);
    const file = await fileHandle.getFile();
    const contents = await file.text();
    let $editor = document.getElementById("editor");
    editor.session.setValue(contents);
    $editor.style.visibility = 'visible';
    $editor.style.display = 'block';
    editor.resize();
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