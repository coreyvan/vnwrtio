// Get the container element
var navContainer = document.getElementById("nav-items");

// Get all buttons with class="btn" inside the container
var items = navContainer.getElementsByClassName("nav-item");

// Loop through the buttons and add the active class to the current/clicked button
for (var i = 0; i < items.length; i++) {
    items[i].addEventListener("click", function () {
        var current = document.getElementsByClassName("active");
        current[0].className = current[0].className.replace(" active under-highlight", "");
        this.className += " active under-highlight";
    });
}

// var words = ["dogs", "security", "code", "Go", "Docker", "jiu jitsu", "networking", "hiking", "swimming", "AWS", "agile", "terraform", "metal working"];
// var speed = 50; /* The speed/duration of the effect in milliseconds */

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function randomIntFromInterval(min, max) { // min and max included 
    return Math.floor(Math.random() * (max - min + 1) + min);
}

const TypingText = document.querySelector(".typed-text");
const AutoTyping = document.querySelector(".TypeCursor");

const textArray = ["stoicism.", "travel.", "dogs.", "microservices.", "security.", "code.", "Go.", "Docker.", "jiu jitsu.", "networking.", "hiking.", "swimming.", "AWS.", "agile.", "terraform.", "cryptography."];
const typingDelay = 100;
const erasingDelay = 80;
const newTextDelay = 1500;
let textArrayIndex = randomIntFromInterval(0, textArray.length);
let charIndex = 0;

function type() {
    if (charIndex < textArray[textArrayIndex].length) {
        if (!AutoTyping.classList.contains("typing")) AutoTyping.classList.add("typing");
        TypingText.textContent += textArray[textArrayIndex].charAt(charIndex);
        charIndex++;
        setTimeout(type, typingDelay);
    }
    else {
        AutoTyping.classList.remove("typing");
        setTimeout(erase, newTextDelay);

    }
}

function erase() {
    if (charIndex > 0) {
        if (!AutoTyping.classList.contains("typing")) AutoTyping.classList.add("typing");
        TypingText.textContent = textArray[textArrayIndex].substring(0, charIndex - 1);
        charIndex--;
        setTimeout(erase, erasingDelay);
    }
    else {
        AutoTyping.classList.remove("typing");
        textArrayIndex = randomIntFromInterval(0, textArray.length)
        if (textArrayIndex >= textArray.length) textArrayIndex = 0;
        setTimeout(type, typingDelay + 1100);
    }
}

document.addEventListener("DOMContentLoaded", function () {
    if (textArray.length) setTimeout(type, newTextDelay + 250);
});