import { useState, useEffect } from "react";

interface TypewriterProps {
	words: string[];
	typingSpeed?: number;
	deletingSpeed?: number;

	delayBetweenWords?: number;
}

export function Typewriter({
	words,

	typingSpeed = 100,

	deletingSpeed = 50,
	delayBetweenWords = 1500,
}: TypewriterProps) {
	const [currentWordIndex, setCurrentWordIndex] = useState(0);
	const [currentText, setCurrentText] = useState("");
	const [isDeleting, setIsDeleting] = useState(false);

	useEffect(() => {
		const word = words[currentWordIndex];

		const timeout = setTimeout(
			() => {
				// If deleting, remove the last character

				if (isDeleting) {
					setCurrentText((prev) => prev.substring(0, prev.length - 1));

					// If all characters are deleted, start typing the next word
					if (currentText === "") {
						setIsDeleting(false);
						setCurrentWordIndex((prev) => (prev + 1) % words.length);
					}
				}
				// If typing, add the next character
				else {
					setCurrentText(word.substring(0, currentText.length + 1));

					// If the word is complete, start deleting after a delay
					if (currentText === word) {
						setTimeout(() => {
							setIsDeleting(true);
						}, delayBetweenWords);
					}
				}
			},
			isDeleting ? deletingSpeed : typingSpeed,
		);

		return () => clearTimeout(timeout);
	}, [
		currentText,
		currentWordIndex,
		isDeleting,
		words,
		typingSpeed,
		deletingSpeed,
		delayBetweenWords,
	]);

	return <span className="text-primary whitespace-nowrap">{currentText}</span>;
}
