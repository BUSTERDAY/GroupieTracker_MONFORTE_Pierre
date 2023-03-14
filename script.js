// Effectuez une requête GET pour les résultats de la recherche
fetch('/search')
.then(response => response.json())
.then(data => {
    // Créez une carte pour chaque résultat
    const container = document.querySelector('.container_cards');
    data.forEach(cocktail => {
        const card = document.createElement('div');
        card.classList.add('card');

        const image = document.createElement('img');
        image.src = cocktail.strDrinkThumb;
        image.alt = cocktail.strDrink;
        card.appendChild(image);

        const name = document.createElement('h2');
        name.textContent = cocktail.strDrink;
        card.appendChild(name);

        const ingredients = document.createElement('ul');
        for (let i = 1; i <= 15; i++) {
            if (cocktail[`strIngredient${i}`]) {
                const li = document.createElement('li');
                li.textContent = `${cocktail[`strIngredient${i}`]} (${cocktail[`strMeasure${i}`]})`;
                ingredients.appendChild(li);
            } else {
                break;
            }
        }
        card.appendChild(ingredients);

        const instructions = document.createElement('p');
        instructions.textContent = cocktail.strInstructions;
        card.appendChild(instructions);
        container.appendChild(card);
});
})
.catch(error => console.error(error));