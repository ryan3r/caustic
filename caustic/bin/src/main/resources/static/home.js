var attempt = 3; // Variable to count number of attempts.
// Below function Executes on click of login button.
function validate() {
    var usrname = document.getElementById("username").value;
    var pssword = document.getElementById("password").value;
    let letNum = new RegExp("^[0-9a-zA-Z]+$");
    var login = false;
    if (usrname.length < 16 && letNum.test(usrname) && pssword.length > 5 && letNum.test(pssword)) {
        var obj = {username: usrname, password: pssword};
        var jsonData = JSON.stringify(obj);
        $.post(URL, jsonData, function(data, status){login = status});
        if(login)
        {
            alert("Login successfully");
            window.location = "formUpload.html";
        }
        return false;
    } else {
        attempt--;// Decrementing by one.
        alert("You have left " + attempt + " attempt;");
// Disabling fields after 3 attempts.
        if (attempt == 0) {
            document.getElementById("username").disabled = true;
            document.getElementById("password").disabled = true;
            document.getElementById("submit").disabled = true;
            return false;
        }
    }
}

function create() {
    var usrname = document.getElementById("username2").value;
    var pssword = document.getElementById("password2").value;
    var act = document.getElementById("actType").value;
    var pswdVer = document.getElementById("pswdVerify").value;
    let letNum = new RegExp("^[0-9a-zA-Z]+$");
    const URL = 'localhost';

    if (usrname.length < 16 && letNum.test(usrname) && pssword.length > 5 && letNum.test(pssword)) {
        alert("Account Created");
        var obj = { username: usrname, password: pssword, actType: act};
        var jsonData = JSON.stringify(obj);
        $.post(URL, jsonData, function(data, status){console.log(`${data} and status is ${status}`)});
        window.location = "formUpload.html";
        return false;
    } else {
        alert("Account not created. Form not valid");
    }
}