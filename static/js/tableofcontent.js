var ToC =
  "<nav role='navigation' class='table-of-contents'>" +
    "<h3>Table of Content</h3>" +
    "<ol class='indexlist' >";

var newLine, el, title, link;

$(".article h2").each(function() {

  el = $(this);
  title = el.text();
  if (title.length > 20) {
  title = title.substring(0, 24);
  }
  setatrr = el.attr("id", title)
  
  link = "#" + title;

  newLine =
    "<li>" +
      "<a href='" + link + "'>" +
        el.text() +
      "</a>" +
    "</li>";

  ToC += newLine;
  

});

ToC +=
   "</ol>" +
  "</nav>";

$(ToC).insertBefore( ".article" )