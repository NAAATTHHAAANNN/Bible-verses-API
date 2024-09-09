function getRandomVerse() {
    fetch('http://localhost:8080/verse')
        .then(response => response.json())
        .then(data => {
            if (data && data.text && data.reference) {
                const verseText = `${data.text} (${data.reference})`;
                document.getElementById('verse').innerText = verseText;
            } else {
                document.getElementById('verse').innerText = "No verse available";
            }
        })
        .catch(error => {
            console.error('Error fetching verse:', error);
            document.getElementById('verse').innerText = "Failed to load verse (check if API is active)";
        });
}

getRandomVerse();