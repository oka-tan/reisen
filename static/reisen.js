//Fix the theme
const linkElement = document.getElementById('main-css-link');
if (window.localStorage.getItem('themeUrl') !== null) {
	const themeUrl = window.localStorage.getItem('themeUrl');
	linkElement.href = themeUrl;
}

//Load event. Most of the code is in here
window.addEventListener('load', function(event) {
	//Global variables
	const quoteLinkRegex = /^>>\d+$/;
	const toggleRegex = /^toggle\('.*'\)$/;
	const board = document.documentElement.getAttribute('data-board');
	const oekakiUrl = document.documentElement.getAttribute('data-oekaki-url');
	const enableLatex = document.documentElement.getAttribute('data-enable-latex') === 'true';
	const enableTegaki = document.documentElement.getAttribute('data-enable-tegaki') === 'true';
	const themeSelect = document.getElementById('theme-select');

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
				if (enableTegaki) {
					const oekakiInternalHash = reisenPost.getAttribute('data-oekaki-internal-hash');
					//Point the href to the post element for convenience
					aElement.href = '#' + reisenPost.id;

					if (oekakiInternalHash) {
						//We need to do this song and dance
						//to inject the hash because of the for loop
						aElement.onclick = function(oekakiInternalHash) {
							return function() {
								Tegaki.open({
									replayMode: true,
									replayURL: `${oekakiUrl}/${oekakiInternalHash}`
								});
							}
						}(oekakiInternalHash);
					} else {
						aElement.textContent += ' (Unavailable)';
					}
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
						return function() {
							const element = document.getElementById(toggleId);

							if (element) {
								element.classList.toggle("hidden");
							}
						}
					}(toggleId);
				}
			}
		}
	}

	//Hide the exif tables by default
	for (const exifTable of document.getElementsByClassName("exif")) {
		exifTable.classList.add("hidden");
	}

	//Mark the correct theme in the select as selected
	if (window.localStorage.getItem('themeName') !== null) {
		const themeName = window.localStorage.getItem('themeName');

		for (const optionElement of themeSelect.getElementsByTagName('option')) {
			optionElement.selected = optionElement.textContent == themeName;

			//Patch the href link in case it's been modified
			if (optionElement.selected && linkElement.href != optionElement.value) {
				window.localStorage.setItem('themeUrl', optionElement.value);
				linkElement.href = optionElement.value;
			}
		}
	}

	themeSelect.onchange = function(event) {
		for (const optionElement of themeSelect.getElementsByTagName('option')) {
			if (optionElement.selected) {
				window.localStorage.setItem('themeUrl', optionElement.value);
				window.localStorage.setItem('themeName', optionElement.textContent);
				linkElement.href = optionElement.value;

				break;
			}
		}
	}


	//Mathjax configuration
	if (enableLatex) {
		console.log('Enabling latex');
		MathJax.Hub.Config({
			tex2jax: {inlineMath: [['[math]','[/math]'], ['[eqn]','[/eqn]']]}
		});
	}
});

