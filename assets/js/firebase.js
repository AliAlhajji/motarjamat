import { initializeApp } from "https://www.gstatic.com/firebasejs/10.7.2/firebase-app.js";
import { getAuth, signInWithEmailAndPassword } from "https://www.gstatic.com/firebasejs/10.7.2/firebase-auth.js";


const firebaseConfig = {

    apiKey: "AIzaSyAO0-1UJjq5Z-BnLDZomAzhtWYoaD9h608",

    authDomain: "motarjamat-2024.firebaseapp.com",

    projectId: "motarjamat-2024",

    storageBucket: "motarjamat-2024.appspot.com",

    messagingSenderId: "795925320860",

    appId: "1:795925320860:web:53ed4a18c28168f6fcc959",

    measurementId: "G-66TV878G6T"

};


// Initialize Firebase

const app = initializeApp(firebaseConfig);

const auth = getAuth(app);

function setCookie(name, value, days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "") + expires + "; path=/";
}


signInWithEmailAndPassword(auth, "al.11.13@hotmail.com", "Ali@123").then((user) => {
    console.log(auth.currentUser);

})

