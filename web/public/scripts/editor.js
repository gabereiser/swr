
function popup_close() {
    let popup = document.getElementById("popup");
    popup.style.visibility = 'hidden';
    popup.getElementsByClassName("popup-content")[0].innerHTML="";
    return false;
}
function menu_about() {
    let popup = document.getElementById("popup");
    popup.style.visibility = 'visible';
    popup.getElementsByClassName("popup-content")[0].innerHTML="About<br/><br/>SWR Editor<br/>&copy; 2022; All rights reserved.";
    return false;
}

function new_file() {
    let $editor = document.getElementById("editor");
    $editor.style.visibility = 'visible';
    editor.session.setValue("");
    editor.resize();
}

async function open_file() {
    [fileHandle] = await window.showOpenFilePicker();
    const file = await fileHandle.getFile();
    const contents = await file.text();
    let $editor = document.getElementById("editor");
    editor.session.setValue(contents);
    $editor.style.visibility = 'visible';
    editor.resize();
    
}

let editor = ace.edit("editor");
    editor.setOption("useSoftTabs", true);
    editor.setTheme("ace/theme/terminal");
    editor.setShowPrintMargin(false);
    editor.session.setMode("ace/mode/yaml");
    editor.resize();