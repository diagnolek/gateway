import React from 'react';
import {NavLink} from "react-router-dom";

const Menu = () => (
    <ul className="nav">
        <li className="nav-item">
            <NavLink className="nav-link" to="/">Home</NavLink>
        </li>
        <li className="nav-item active">
            <a href="#">Network</a>
            <ul className="nav">
                <li className="nav-item">
                    <NavLink className="nav-link" to="/network/ethernet">&gt;ethernet</NavLink>
                </li>
            </ul>
        </li>
        <li className="nav-item">
            <a href="#">Feature</a>
            <ul className="nav">
                <li className="nav-item">
                    <NavLink className="nav-link" to="/feature/chat">&gt;chat user</NavLink>
                </li>
            </ul>
        </li>
    </ul>
);

export default Menu;