//Global functions
function playOekaki(oekakiInternalHash) {
	Tegaki.open({
		replayMode: true,
		replayURL: `http://s3.localhost/ayase/${oekakiInternalHash}`
	});
}

//Global variables
var quoteLinkRegex = /^>>\d+$/;

//Doing stuff proper
for (let reisenPost of document.getElementsByClassName('reisen-post')) {
	postNumber = reisenPost.id.substring(1);

	for (let aElement of reisenPost.getElementsByTagName('a')) {
		if (aElement.textContent.match(quoteLinkRegex)) {
			postNumberReferenced = aElement.textContent.substring(2);

			quoteLinksHolder = document.getElementById('quoteLinks' + postNumberReferenced);

			if (quoteLinksHolder) {
				quoteLink = document.createElement('a');

				quoteLink.textContent = '>>' + postNumber;
				quoteLink.href = '#' + reisenPost.id;
				quoteLink.classList.add('text-violet-500');
				quoteLink.classList.add('hover:text-violet-600');
				quoteLink.classList.add('text-xs');

				quoteLinksHolder.appendChild(quoteLink);
			}
		} else if (aElement.href.substring(0, 19) === 'javascript:oeReplay') {
			oekakiInternalHash = reisenPost.getAttribute('data-oekaki-internal-hash');
			aElement.href = '#' + reisenPost.id;

			if (oekakiInternalHash) {
				aElement.onclick = function(oekakiInternalHash) {
					return function() { playOekaki(oekakiInternalHash); }
				}(oekakiInternalHash)
			} else {
				aElement.textContent += ' (Unavailable)';
			}
		}
	}
}

MathJax.Hub.Config({
	tex2jax: {inlineMath: [['[math]','[/math]'], ['[eqn]','[/eqn]']]}
});

