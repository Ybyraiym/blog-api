const createForm = document.getElementById('create-form');
const postList = document.getElementById('post-list');

createForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    
    const title = document.getElementById('title').value;
    const content = document.getElementById('content').value;
    
    const response = await fetch('/blogposts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ title, content })
    });
    
    if (response.ok) {
        const newPost = await response.json();
        addPostToList(newPost);
        createForm.reset();
    }
});

async function fetchPosts() {
    const response = await fetch('/blogposts');
    const posts = await response.json();
    posts.forEach(addPostToList);
}

function addPostToList(post) {
    const listItem = document.createElement('li');
    listItem.innerHTML = `<strong>${post.title}</strong><br>${post.content}`;
    postList.appendChild(listItem);
}

fetchPosts();