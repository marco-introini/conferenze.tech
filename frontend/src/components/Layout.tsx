import { type ReactNode } from "react";
import { Link } from "react-router-dom";
import { Map } from "lucide-react";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen flex flex-col bg-slate-50 font-sans">
      <nav className="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200/50">
        <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-3">
            <div className="w-9 h-9 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-indigo-500/25">
              <Map className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              conferenze.tech
            </span>
          </Link>
          <div className="flex items-center gap-3">
            <Link to="/login" className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-slate-900 transition-colors">
              Accedi
            </Link>
            <Link to="/register" className="px-5 py-2 text-sm font-semibold bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-xl hover:shadow-lg hover:shadow-indigo-500/25 transition-all duration-300">
              Registrati
            </Link>
          </div>
        </div>
      </nav>

      <div className="flex-1 pt-16">
        {children}
      </div>

      <footer className="border-t border-slate-200 bg-white mt-auto">
        <div className="max-w-6xl mx-auto px-6 py-8 flex items-center justify-between">
          <p className="text-slate-500 text-sm">© 2024 conferenze.tech</p>
          <div className="flex gap-6">
            <a href="#" className="text-slate-500 hover:text-slate-900 text-sm transition-colors">Privacy</a>
            <a href="#" className="text-slate-500 hover:text-slate-900 text-sm transition-colors">Terms</a>
            <a href="#" className="text-slate-500 hover:text-slate-900 text-sm transition-colors">Contact</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
