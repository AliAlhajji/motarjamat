import { getAuth, signInWithEmailAndPassword, signOut, onAuthStateChanged } from "https://www.gstatic.com/firebasejs/10.7.2/firebase-auth.js";
import { setCookie } from "./utils.js";

const auth = getAuth()


onAuthStateChanged(auth, (user) => {

    auth.currentUser.getIdToken().then((token) => {
        setCookie("token", token, 1)
    }).catch((error) => {
        console.log(error)
    })


})

export function signIn() {
    var email = document.getElementById("email").value
    var password = document.getElementById("password").value
    var loginError = document.getElementById("loginError")

    if (email == "" || password == "") {
        loginError.textContent = "جميع الحقول مطلوبة"
        loginError.hidden = false
        return
    }
    signInWithEmailAndPassword(auth, email, password).then((creds) => {

        if (creds.user) {
            setCookie("token", creds.user.accessToken, 1)
            window.location = "/home"
            return
        }
        else {
            loginError.textContent = "معلومات الدخول غير صحيحة"
            loginError.hidden = false
            return
        }
    }).catch((error) => {
        if (error.code == "auth/invalid-credential") {
            loginError.textContent = "معلومات الدخول غير صحيحة"
            loginError.hidden = false
        }
        else {
            loginError.textContent = error
            loginError.hidden = false
            return
        }
    });

}

export function register() {
    var email = document.getElementById("email").value
    var password = document.getElementById("password").value
    var username = document.getElementById("username").value
    var name = document.getElementById("name").value

    var form = document.getElementById("form")

    var errorMsg = document.getElementById("errorMsg")

    if (email == "" || password == "" || username == "" || name == "") {
        errorMsg.textContent = "كل الحقول مطلوبة"
        errorMsg.hidden = false
        return
    }

    form.submit()


}

export function logout() {

    signOut(auth).then((u) => {
        setCookie("token", "", 0)
        location.reload()
    })
}