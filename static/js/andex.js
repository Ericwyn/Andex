function gotoDir(path) {
    if (path === null || path === undefined){
        return
    }
    console.log("goto path: " + path)
    if (path === "/"){
        window.location.href= window.location.origin
    } else {
        window.location.search =  "?p=" + path
    }
}

function downFile() {
    let filePath = window.location.search;
    filePath = filePath.replace("/?p=", "")
    filePath = filePath.replace("?p=", "")
    let newUrl = "download?p=" + filePath
    window.open(newUrl,"_blank")
}