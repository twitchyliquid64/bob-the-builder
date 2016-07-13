function flashHighlightElement(element){
	element.addClass("flashHighlight");
  setTimeout( function(){
        element.removeClass("flashHighlight");
    }, 1500);	// Timeout must be the same length as the CSS3 transition or longer (or you'll mess up the transition)
}
