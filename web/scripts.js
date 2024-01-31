document.addEventListener('DOMContentLoaded', function () {
    const networkForm = document.getElementById('networkForm');
    const networkList = document.getElementById('networkList');

    networkForm.addEventListener('submit', function (event) {
        event.preventDefault();

        const name = document.getElementById('name').value;
        const subnet = document.getElementById('subnet').value;
        const parent_id = document.getElementById('parent_id').value;

        fetch('/api/networks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, subnet, parent_id }),
        })
        .then(response => response.json())
        .then(data => {
            // Clear form inputs
            document.getElementById('name').value = '';
            document.getElementById('subnet').value = '';
            document.getElementById('parent_id').value = '';

            // Display the newly created network in the list
            const networkItem = document.createElement('div');
            networkItem.textContent = `ID: ${data.id}, Subnet: ${data.subnet}, Parent ID: ${data.parent_id || 'None'}`;
            networkList.appendChild(networkItem);
        })
        .catch(error => console.error('Error:', error));
    });

    // Fetch and display existing networks on page load
    fetch('/api/networks')
    .then(response => response.json())
    .then(data => {
        data.forEach(network => {
            const networkItem = document.createElement('div');
            networkItem.textContent = `ID: ${network.id}, Subnet: ${network.subnet}, Parent ID: ${network.parent_id || 'None'}`;
            networkList.appendChild(networkItem);
        });
    })
    .catch(error => console.error('Error:', error));
});