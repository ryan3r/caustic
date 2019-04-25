function checkFile(e)
{
    var file = document.getElementById('upload').value;
    var extension = file.substr((file.lastIndexOf('.') +1));
    if (!/(cpp|py|java)$/ig.test(extension)) {
        alert("Invalid file type: "+extension+".  File must be Java, C++, or Python.");
        $("#file").val("");
        e.preventDefault();
    }

}

document.querySelector("form").addEventListener("submit", checkFile);