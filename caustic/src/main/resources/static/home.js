var attempt = 3; // Variable to count number of attempts.
// Below function Executes on click of login button.
function validate() {
    var usrname = document.getElementById("username").value;
    var pssword = document.getElementById("password").value;
    let letNum = new RegExp("^[0-9a-zA-Z]+$");
    var login = false;
    var tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const URL = '/accountsLogin';
    if (usrname.length < 16 && letNum.test(usrname) && pssword.length > 5 && letNum.test(pssword)) {
        var obj = {username: usrname, password: pssword};
        var jsonData = JSON.stringify(obj);
        var xmlhttp = new XMLHttpRequest();
        xmlhttp.open("POST", URL, true);
        xmlhttp.setRequestHeader("Content-Type", "application/json");
        xmlhttp.onreadystatechange = function () {
            if(xmlhttp.readyState === 4 && xmlhttp.status === 200) {
              login = JSON.parse(xmlhttp.responseText);
                if(!login)
                {
                    alert("Login failed");
                }
                else
                {
                    alert("Login successfully");
                    document.cookie = "username=" + usrname + "; expires=" + tomorrow + "; path=/";
                    window.location.href = "formUpload";
                }
            }
        };
        xmlhttp.send(JSON.stringify(obj));
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
    var acc = document.getElementById("actType").value;
    var pswdVer = document.getElementById("pswdVerify").value;
    let letNum = new RegExp("^[0-9a-zA-Z]+$");
    const URL = '/accounts';
    var tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    if (usrname.length < 16 && letNum.test(usrname) && pssword.length > 5 && letNum.test(pssword)) {
        alert("Account Created");
        var obj = { username: usrname, password: pssword, accType: acc};
        var xmlhttp = new XMLHttpRequest();
        var bool = false;
        xmlhttp.open("POST", URL, true);
        xmlhttp.setRequestHeader("Content-Type", "application/json");
        xmlhttp.onreadystatechange = function () {
            if(xmlhttp.readyState === 4 && xmlhttp.status === 200) {
              bool = JSON.parse(xmlhttp.responseText);
		    if(!bool)
		    {
	            alert("Username already in use");
        	}
	        else
    		{
	            document.cookie = "username=" + usrname + "; expires=" + tomorrow + "; path=/";
	            window.location.href = "formUpload";
        	}
        }
    };
	xmlhttp.send(JSON.stringify(obj));
        
        return false;
    } else {
        alert("Account not created. Form not valid");
    }
}