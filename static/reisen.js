//Global variables
const quoteLinkRegex = /^>>\d+$/;
const toggleRegex = /^toggle\('.*'\)$/;
const board = document.documentElement.getAttribute('data-board');
const oekakiUrl = document.documentElement.getAttribute('data-oekaki-url');

//Global functions
function playOekaki(oekakiInternalHash) {
	Tegaki.open({
		replayMode: true,
		replayURL: `${oekakiUrl}/${oekakiInternalHash}`
	});
}

function toggle(id) {
	const element = document.getElementById(id);

	if (element) {
		element.classList.toggle("hidden");
	}
}

//Doing stuff proper
for (const reisenPost of document.getElementsByClassName('reisen-post')) {
	const postNumber = reisenPost.id.substring(1);


	for (const deadLinkElement of reisenPost.getElementsByClassName('deadlink')) {
		const aElement = document.createElement('a');
		
		aElement.textContent = deadLinkElement.textContent;
		
		deadLinkElement.parentNode.replaceChild(aElement, deadLinkElement);
	}

	for (const aElement of reisenPost.getElementsByTagName('a')) {
		if (aElement.classList.contains('reisen-backlink')) {
			continue;
		} else if (aElement.textContent.match(quoteLinkRegex)) {
			const postNumberReferenced = aElement.textContent.substring(2);

			//Fix search quotelinks
			if (reisenPost.classList.contains('reisen-post-ambiguous') || reisenPost.classList.contains('reisen-thread')) {
				aElement.href = '/' + board + '/post/' + postNumberReferenced;
			} else if (reisenPost.classList.contains('reisen-post-op') || reisenPost.classList.contains('reisen-post-reply')) {
				const quoteLinksHolder = document.getElementById('quoteLinks' + postNumberReferenced);

				if (quoteLinksHolder) {
					//Ensure the link is correct.
					aElement.href = '#p' + postNumberReferenced;

					const quoteLink = document.createElement('a');

					quoteLink.textContent = '>>' + postNumber;
					quoteLink.href = '#' + reisenPost.id;
					quoteLink.classList.add('reisen-backlink');

					quoteLinksHolder.appendChild(quoteLink);
				} else {
					//Correct the link.
					aElement.href = '/' + board + '/post/' + postNumberReferenced;
				}
			}
		} else if (aElement.href.substring(0, 19) === 'javascript:oeReplay') {
			//Fix oekaki links
			const oekakiInternalHash = reisenPost.getAttribute('data-oekaki-internal-hash');
			//Point the href to the post element for convenience
			aElement.href = '#' + reisenPost.id;

			if (oekakiInternalHash) {
				//We need to do this song and dance
				//to inject the hash because of the for loop
				aElement.onclick = function(oekakiInternalHash) {
					return function() { playOekaki(oekakiInternalHash); }
				}(oekakiInternalHash)
			} else {
				aElement.textContent += ' (Unavailable)';
			}
		} else if (aElement.getAttribute('onclick')) {
			
			//The href is something like 'javascript:void(0)',
			//which causes CSP issues with the onclick action.
			aElement.href = '#' + reisenPost.id;

			const rawOnclick = aElement.getAttribute('onclick');
			//CSP should, unless you've misconfigured it,
			//prevent the onclick from working (which is good and secure).
			//Hence, we need to fix it (but just toggles).
			if (rawOnclick.match(toggleRegex)) {
				const toggleId = rawOnclick.substring(8, rawOnclick.length - 2);
				aElement.onclick = function(toggleId) {
					return function() { toggle(toggleId); }
				}(toggleId);
			}
		}
	}
}

//Hide the exif tables by default
for (const exifTable of document.getElementsByClassName("exif")) {
	exifTable.classList.add("hidden");
}

//Mathjax configuration
MathJax.Hub.Config({
	tex2jax: {inlineMath: [['[math]','[/math]'], ['[eqn]','[/eqn]']]}
});

