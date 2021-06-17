
    window.addEventListener("load", async () => {
    let categoryChevron = document.querySelectorAll(".category_arrow")
    categoryChevron.forEach((e) => {
        e.addEventListener("click", () =>{
            let display = e.parentNode.parentNode.querySelector(".contents").style.display
            if(display == "none"){
                e.classList.add("categorie_chevron_active")
                e.parentNode.parentNode.querySelector(".contents").style.display = "flex"
            }else{
                e.classList.remove("categorie_chevron_active")
                e.parentNode.parentNode.querySelector(".contents").style.display = "none"
            }
        });
    })
})
let categoryChevron = document.querySelectorAll(".category_arrow")
categoryChevron.forEach((e) => {
    e.addEventListener("click", () =>{
        let display = e.parentNode.parentNode.querySelector(".contents").style.display
        if(display == "none"){
            e.classList.add("categorie_chevron_active")
            e.parentNode.parentNode.querySelector(".contents").style.display = "flex"
        }else{
            e.classList.remove("categorie_chevron_active")
            e.parentNode.parentNode.querySelector(".contents").style.display = "none"
        }
    });
});

let navChevron = document.getElementById("nav-chevron");
let navDropdown = document.getElementById("nav-dropdown");
navChevron.addEventListener("mouseover", () => {
    navDropdown.classList.toggle("active");
});


if (document.body.innerHTML.includes("sharing")) {
    document.querySelector('.sharing').addEventListener('click', function() {
    copy()
    });
};

function copy(){
let inputElement = document.createElement("input");
inputElement.type = "text";
inputElement.value = window.location.href;
document.body.appendChild(inputElement);
inputElement.select();
document.execCommand('copy');
document.body.removeChild(inputElement);
window.alert("Post's link copied!")
}
    all_img = document.querySelectorAll('img')
    all_img.forEach(function(img){
        if (img.getAttribute('src') == "") {
        img.remove()
        }
    });
