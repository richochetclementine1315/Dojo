import { Link } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { Button } from '../ui/Button';
import { Menu, X, User, LogOut, Settings } from 'lucide-react';
import { useState } from 'react';

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  const closeMobileMenu = () => {
    setIsMobileMenuOpen(false);
  };

  return (
    <nav className="sticky top-0 z-50 w-full border-b border-dojo-red-900/20 bg-dojo-black-900/95 backdrop-blur">
      <div className="container mx-auto px-4 flex h-16 items-center justify-between">
        <Link to="/" className="cursor-target flex items-center group">
          <img 
            src="/dojo-logo.png" 
            alt="Dojo Logo" 
            className="h-16 w-16 object-contain group-hover:scale-110 transition-transform" 
          />
        </Link>

        {/* Desktop Nav Links */}
        {isAuthenticated && (
          <div className="hidden md:flex items-center space-x-6">
            <Link 
              to="/dashboard" 
              className="cursor-target text-gray-300 hover:text-white transition-colors font-medium"
            >
              Dashboard
            </Link>
            <Link 
              to="/problems" 
              className="cursor-target text-gray-300 hover:text-white transition-colors font-medium"
            >
              Problems
            </Link>
            <Link 
              to="/contests" 
              className="cursor-target text-gray-300 hover:text-white transition-colors font-medium"
            >
              Contests
            </Link>
            <Link 
              to="/sheets" 
              className="cursor-target text-gray-300 hover:text-white transition-colors font-medium"
            >
              Sheets
            </Link>
            <Link 
              to="/rooms" 
              className="cursor-target text-gray-300 hover:text-white transition-colors font-medium"
            >
              Rooms
            </Link>
            <Link 
              to="/profile" 
              className="cursor-target text-gray-300 hover:text-white transition-colors font-medium"
            >
              Profile Stats
            </Link>
          </div>
        )}

        {/* Desktop Auth Section */}
        <div className="hidden md:flex items-center space-x-4">
          {isAuthenticated ? (
            <div className="relative">
              <button
                onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                className="cursor-target flex items-center space-x-2 px-3 py-2 rounded-lg hover:bg-dojo-black-800 transition-colors"
              >
                <div className="h-8 w-8 rounded-full bg-gradient-to-r from-dojo-red-500 to-orange-500 flex items-center justify-center text-white font-semibold">
                  {user?.username?.charAt(0).toUpperCase() || 'U'}
                </div>
                <span className="text-gray-300">{user?.username}</span>
              </button>

              {isUserMenuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-dojo-black-800 border border-gray-700 rounded-lg shadow-lg py-2">
                  <Link
                    to="/profile"
                    onClick={() => setIsUserMenuOpen(false)}
                    className="cursor-target flex items-center space-x-2 px-4 py-2 text-gray-300 hover:bg-dojo-black-700 hover:text-white transition-colors"
                  >
                    <User className="h-4 w-4" />
                    <span>Profile</span>
                  </Link>
                  <Link
                    to="/settings/platforms"
                    onClick={() => setIsUserMenuOpen(false)}
                    className="cursor-target flex items-center space-x-2 px-4 py-2 text-gray-300 hover:bg-dojo-black-700 hover:text-white transition-colors"
                  >
                    <Settings className="h-4 w-4" />
                    <span>Platform Settings</span>
                  </Link>
                  <button
                    onClick={() => {
                      setIsUserMenuOpen(false);
                      handleLogout();
                    }}
                    className="cursor-target w-full flex items-center space-x-2 px-4 py-2 text-gray-300 hover:bg-dojo-black-700 hover:text-white transition-colors"
                  >
                    <LogOut className="h-4 w-4" />
                    <span>Logout</span>
                  </button>
                </div>
              )}
            </div>
          ) : (
            <>
              <Link to="/login">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link to="/register">
                <Button>Sign Up</Button>
              </Link>
            </>
          )}
        </div>

        {/* Mobile Menu Button */}
        <button
          className="cursor-target md:hidden text-white"
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
        >
          {isMobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
        </button>
      </div>

      {/* Mobile Menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-gray-700 bg-dojo-black-900">
          <div className="container mx-auto px-4 py-4 space-y-3">
            {isAuthenticated ? (
              <>
                <div className="flex items-center space-x-3 pb-3 border-b border-gray-700">
                  <div className="h-10 w-10 rounded-full bg-gradient-to-r from-dojo-red-500 to-orange-500 flex items-center justify-center text-white font-semibold">
                    {user?.username?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <span className="text-white font-medium">{user?.username}</span>
                </div>
                <Link 
                  to="/dashboard" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Dashboard
                </Link>
                <Link 
                  to="/problems" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Problems
                </Link>
                <Link 
                  to="/contests" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Contests
                </Link>
                <Link 
                  to="/sheets" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Sheets
                </Link>
                <Link 
                  to="/rooms" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Rooms
                </Link>
                <Link 
                  to="/profile" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Profile Stats
                </Link>
                <Link 
                  to="/settings/platforms" 
                  onClick={closeMobileMenu}
                  className="cursor-target block py-2 text-gray-300 hover:text-white transition-colors"
                >
                  Platform Settings
                </Link>
                <button
                  onClick={() => {
                    closeMobileMenu();
                    handleLogout();
                  }}
                  className="cursor-target w-full text-left py-2 text-red-400 hover:text-red-300 transition-colors"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" onClick={closeMobileMenu}>
                  <Button variant="ghost" className="w-full justify-start">Login</Button>
                </Link>
                <Link to="/register" onClick={closeMobileMenu}>
                  <Button className="w-full">Sign Up</Button>
                </Link>
              </>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}
