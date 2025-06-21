import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Header from "../../components/Header/Header";
import Menu from "../../components/Menu/Menu";
import HomeView from "../HomeView/HomeView";
import EthernetView from "../NetworkView/EthernetView";
import ChatView from "../FeatureView/ChatView";
import Logo from "../../components/Header/Logo";

const Root = () => {

    return (
        <BrowserRouter>
            <div className="container">
                <div className="columns root-container">
                    <div className="column col-3 root-menu">
                        <Logo />
                        <Menu />
                    </div>
                    <div className="column col-9 root-content">
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