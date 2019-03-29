var attempt = 3; // Variable to count number of attempts.
// Below function Executes on click of login button.
function validate() {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    let letNum = new RegExp("^[0-9a-zA-Z]+$");

    if (username.length < 16 && letNum.test(username) && password.length > 5 && letNum.test(password)) {
        alert("Login successfully");
        window.location = "formUpload.html";
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
    var username = document.getElementById("username2").value;
    var password = document.getElementById("password2").value;
    var pswdVer = document.getElementById("pswdVerify").value;
    let letNum = new RegExp("^[0-9a-zA-Z]+$");

    if (username.length < 16 && letNum.test(username) && password.length > 5 && letNum.test(password)) {
        alert("Account Created");
        window.location = "formUpload.html";
        return false;
    } else {
        alert("Account not created. Form not valid");
    }
}