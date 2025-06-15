import React, { useState, useEffect } from 'react';

const Root = () => {
    const [fetchedText, setFetchedText] = useState('');

    useEffect(() => {
        fetch('http://localhost:8080')
            .then((response) => response.text())
            .then((data) => setFetchedText(data))
            .catch((err) => console.error('Error fetching data:', err));
    }, []);

    return (
        <p>This is a Root View, {fetchedText}</p>
    );
};

export default Root;