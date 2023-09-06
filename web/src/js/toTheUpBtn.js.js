document.addEventListener("scroll", handleScroll);
// get a reference to our predefined button
var scrollToTopBtn = document.querySelector(".scrollToTopBtn");

function handleScroll() {
  var scrollableHeight =
    document.documentElement.scrollHeight -
    document.documentElement.clientHeight;
  var GOLDEN_RATIO = 0.3;

  if (document.documentElement.scrollTop / scrollableHeight > GOLDEN_RATIO) {
    //show button
    scrollToTopBtn.classList.remove("hidden");
    // scrollToTopBtn.style.display = "block";
  } else {
    //hide button
    scrollToTopBtn.classList.add("hidden");
    // scrollToTopBtn.style.display = "none";
  }
}

scrollToTopBtn.addEventListener("click", scrollToTop);

function scrollToTop() {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
}
