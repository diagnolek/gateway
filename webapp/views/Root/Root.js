import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Header from "../../components/Header/Header";
import Menu from "../../components/Menu/Menu";
import HomeView from "../HomeView/HomeView";
import EthernetView from "../NetworkView/EthernetView";
import ChatView from "../FeatureView/ChatView";

const Root = () => {
    // const [fetchedText, setFetchedText] = useState('');
    //
    // useEffect(() => {
    //     fetch('http://localhost:8080')
    //         .then((response) => response.text())
    //         .then((data) => setFetchedText(data))
    //         .catch((err) => console.error('Error fetching data:', err));
    // }, []);

    return (
        <BrowserRouter>
            <div className="container">
                <div className="columns">
                    <div className="column col-3">
                        <Menu />
                    </div>
                    <div className="column col-9">
                        <Header />
                        <Routes>
                            <Route path="/" element={<HomeView />} />
                            <Route path="/network/ethernet" element={<EthernetView />} />
                            <Route path="/feature/chat" element={<ChatView />} />
                        </Routes>
                    </div>
                </div>
            </div>
        </BrowserRouter>
    );
};

export default Root;