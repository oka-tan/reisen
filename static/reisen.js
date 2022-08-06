//Global functions
function playOekaki(oekakiInternalHash) {
	Tegaki.open({
		replayMode: true,
		replayURL: `http://s3.localhost/ayase/${oekakiInternalHash}`
	});
}

function toggle(id) {
	var element = document.getElementById(id);
	element.classList.toggle("hidden");
}

//Global variables
var quoteLinkRegex = /^>>\d+$/;
var toggleRegex = /^toggle\('.*'\)$/;
var board = document.documentElement.getAttribute('data-board');

//Doing stuff proper
for (let reisenPost of document.getElementsByClassName('reisen-post')) {
	postNumber = reisenPost.id.substring(1);

	for (let aElement of reisenPost.getElementsByTagName('a')) {
		if (aElement.textContent.match(quoteLinkRegex)) {
			postNumberReferenced = aElement.textContent.substring(2);

			//Fix search quotelinks
			if (reisenPost.classList.contains('reisen-post-ambiguous') || reisenPost.classList.contains('reisen-thread')) {
				aElement.href = '/' + board + '/post/' + postNumberReferenced;
			} else if (reisenPost.classList.contains('reisen-post-op') || reisenPost.classList.contains('reisen-post-reply')) {
				quoteLinksHolder = document.getElementById('quoteLinks' + postNumberReferenced);

				if (quoteLinksHolder) {
					//Ensure the link is correct.
					aElement.href = '#p' + postNumberReferenced;

					quoteLink = document.createElement('a');

					quoteLink.textContent = '>>' + postNumber;
					quoteLink.href = '#' + reisenPost.id;
					quoteLink.classList.add('text-violet-500');
					quoteLink.classList.add('hover:text-violet-600');
					quoteLink.classList.add('text-xs');

					quoteLinksHolder.appendChild(quoteLink);
				} else {
					//Correct the link.
					aElement.href = '/' + board + '/post/' + postNumberReferenced;
				}
			}
		} else if (aElement.href.substring(0, 19) === 'javascript:oeReplay') {
			//Fix oekaki links
			oekakiInternalHash = reisenPost.getAttribute('data-oekaki-internal-hash');
			//Point the href to the post element for convenience
			aElement.href = '#' + reisenPost.id;
			aElement.classList.add('text-violet-500');
			aElement.classList.add('hover:text-violet-600');

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
			aElement.classList.add('text-violet-500');
			aElement.classList.add('hover:text-violet-600');

			rawOnclick = aElement.getAttribute('onclick');
			//CSP should, unless you've misconfigured it,
			//prevent the onclick from working (which is good and secure).
			//Hence, we need to fix it (but just toggles).
			if (rawOnclick.match(toggleRegex)) {
				toggleId = rawOnclick.substring(8, rawOnclick.length - 2);
				aElement.onclick = function(toggleId) {
					return function() { toggle(toggleId); }
				}(toggleId);
			}
		}
	}
}

//Hide the exif tables by default
for (let exifTable of document.getElementsByClassName("exif")) {
	exifTable.classList.add("hidden");
}

//Mathjax configuration
MathJax.Hub.Config({
	tex2jax: {inlineMath: [['[math]','[/math]'], ['[eqn]','[/eqn]']]}
});

