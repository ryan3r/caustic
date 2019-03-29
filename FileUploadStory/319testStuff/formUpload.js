function checkFile()
{
    var file = document.getElementById('upload').value;
    var extension = file.substr((file.lastIndexOf('.') +1));
    if (!/(cpp|py|java)$/ig.test(extension)) {
        alert("Invalid file type: "+extension+".  File must be Java, C++, or Python.");
        $("#file").val("");
    }
    else
    {
        alert("Thank you for the file");
    }

}