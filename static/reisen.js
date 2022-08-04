MathJax.Hub.Config({
	tex2jax: {inlineMath: [['[math]','[/math]'], ['[eqn]','[/eqn]']]}
});

function playOekaki(oekakiInternalHash) {
	Tegaki.open({
		replayMode: true,
		replayURL: `http://s3.localhost/ayase/${oekakiInternalHash}`
	});
}

for (let reisenPost of document.getElementsByClassName("reisen-post")) {
	for (let aElement of reisenPost.getElementsByTagName("a")) {
		if (aElement.href.substring(0, 19) === "javascript:oeReplay") {
			oekakiInternalHash = reisenPost.getAttribute("data-oekaki-internal-hash");
			aElement.href = "#" + reisenPost.id;

			if (oekakiInternalHash) {
				aElement.onclick = function(oekakiInternalHash) {
					return function() { playOekaki(oekakiInternalHash); }
				}(oekakiInternalHash)
			} else {
				aElement.textContent += " (Unavailable)";
			}
		}
	}
}
