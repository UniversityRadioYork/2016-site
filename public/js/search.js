let form = document.getElementById("searchForm");
let filters = document.getElementsByClassName("search-filter")
for (var i = 0; i < filters.length; i++) {
  filters[i].addEventListener("click", function () {
    form.submit();
  });
};
