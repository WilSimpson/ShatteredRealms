import { MemoryRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import Login from './pages/Login';
import DefaultTemplate from './templates/default';
import Home from './pages/Home';

export default function App() {
  return (
    <DefaultTemplate>
      <Router>
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/home" element={<Home />} />
        </Routes>
      </Router>
    </DefaultTemplate>
  );
}
