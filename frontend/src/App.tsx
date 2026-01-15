import { useState, useRef, useEffect } from "react";
import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import ConferenceCard from "./components/ConferenceCard";
import ConferenceMap from "./components/ConferenceMap";
import Layout from "./components/Layout";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Profilo from "./pages/Profilo";
import MieConferenze from "./pages/MieConferenze";
import RegistrazioneConferenza from "./pages/RegistrazioneConferenza";
import type { Conference } from "./types";
import { Map, Calendar, Users, Grid, Search } from "lucide-react";
import { AuthProvider, useAuth } from "./AuthContext";
import { api } from "./api";
import "./index.css";

function Home() {
  const { isAuthenticated, user } = useAuth();
  const [conferences, setConferences] = useState<Conference[]>([]);
  const [viewMode, setViewMode] = useState<"list" | "map">("map");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedConferenceId, setSelectedConferenceId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const selectedRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const loadConferences = async () => {
      try {
        const data = await api.getConferences();
        setConferences(data);
      } catch (err) {
        console.error("Failed to load conferences:", err);
        setError("Impossibile caricare le conferenze");
      } finally {
        setIsLoading(false);
      }
    };
    loadConferences();
  }, []);

  const filteredConferences = conferences.filter(conf => {
    const query = searchQuery.toLowerCase();
    return (
      conf.title.toLowerCase().includes(query) ||
      conf.location.toLowerCase().includes(query)
    );
  });

  const handleSelectConference = (id: string) => {
    setSelectedConferenceId(id);
    setViewMode("list");
    setTimeout(() => {
      selectedRef.current?.scrollIntoView({ behavior: "smooth", block: "center" });
    }, 100);
  };

  const selectedConference = selectedConferenceId
    ? conferences.find(c => c.id === selectedConferenceId)
    : null;

  return (
    <Layout>
      <section className="relative pt-32 pb-20 px-6 overflow-hidden">
        <div className="absolute inset-0 bg-linear-to-br from-indigo-50 via-white to-purple-50" />
        <div className="absolute top-20 right-0 w-96 h-96 bg-linear-to-br from-indigo-200/30 to-purple-200/30 rounded-full blur-3xl" />
        <div className="absolute bottom-0 left-0 w-72 h-72 bg-linear-to-tr from-purple-200/30 to-indigo-200/30 rounded-full blur-3xl" />
        
        <div className="relative max-w-6xl mx-auto">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between">
            <div className="max-w-3xl">
              <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/80 backdrop-blur-sm rounded-full border border-slate-200/50 shadow-sm mb-6">
                <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                <span className="text-sm font-medium text-slate-600">
                  {isLoading ? "Caricamento..." : `${conferences.length} eventi in programma`}
                </span>
              </div>
              <h1 className="text-5xl font-bold text-slate-900 leading-tight mb-6">
                Trova conferenze tech e
                <span className="bg-linear-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent"> condividi il viaggio</span>
              </h1>
              <p className="text-xl text-slate-600 mb-8 max-w-2xl leading-relaxed">
                Connettiti con altri sviluppatori, riduci i costi di viaggio e trova compagnia per il tuo prossimo evento tech.
              </p>
            </div>
          </div>

          <div className="flex flex-wrap gap-4 mt-8">
            <button
              onClick={() => setViewMode("map")}
              className="group px-6 py-3 bg-linear-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl shadow-lg shadow-indigo-500/25 hover:shadow-xl hover:shadow-indigo-500/30 transition-all duration-300 flex items-center gap-2"
            >
              <Map className="w-4 h-4" />
              Esplora sulla mappa
            </button>
            {!isAuthenticated && (
              <Link
                to="/register"
                className="px-6 py-3 bg-white text-slate-700 font-semibold rounded-xl border border-slate-200 hover:border-slate-300 hover:bg-slate-50 transition-all duration-300"
              >
                Registrati
              </Link>
            )}
          </div>
        </div>
      </section>

      <main className="max-w-6xl mx-auto px-6 py-16">
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">
            {error}
          </div>
        )}

        <div className="mb-10">
          <div className="flex flex-col lg:flex-row lg:items-end justify-between gap-6 mb-8">
            <div>
              <h2 className="text-3xl font-bold text-slate-900 mb-2">
                {viewMode === "map" ? "Mappa eventi" : "Tutti gli eventi"}
              </h2>
              <p className="text-slate-500">
                {viewMode === "map"
                  ? "Clicca sui marker per vedere i dettagli"
                  : "Sfoglia le prossime conferenze tech in Italia"}
              </p>
            </div>
            
            <div className="flex flex-col sm:flex-row gap-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <input
                  type="text"
                  placeholder="Cerca per città o evento..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent w-full sm:w-64"
                />
              </div>
              
              <div className="flex bg-white border border-slate-200 rounded-xl p-1">
                <button
                  onClick={() => setViewMode("list")}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 flex items-center gap-2 ${
                    viewMode === "list"
                      ? "bg-indigo-600 text-white"
                      : "text-slate-600 hover:text-slate-900"
                  }`}
                >
                  <Grid className="w-4 h-4" />
                  Lista
                </button>
                <button
                  onClick={() => setViewMode("map")}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 flex items-center gap-2 ${
                    viewMode === "map"
                      ? "bg-indigo-600 text-white"
                      : "text-slate-600 hover:text-slate-900"
                  }`}
                >
                  <Map className="w-4 h-4" />
                  Mappa
                </button>
              </div>
            </div>
          </div>

          {isLoading ? (
            <div className="flex items-center justify-center py-20">
              <span className="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin" />
            </div>
          ) : viewMode === "map" ? (
            <ConferenceMap
              conferences={conferences}
              onSelectConference={handleSelectConference}
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredConferences.map((conf) => (
                <div
                  key={conf.id}
                  ref={conf.id === selectedConferenceId ? selectedRef : null}
                  className={`transition-all duration-300 ${
                    conf.id === selectedConferenceId
                      ? "ring-2 ring-indigo-500 rounded-2xl"
                      : ""
                  }`}
                >
                  <ConferenceCard conf={conf} />
                </div>
              ))}
              <div className="group relative overflow-hidden rounded-2xl border-2 border-dashed border-slate-300 hover:border-indigo-400 transition-colors duration-300 flex items-center justify-center min-h-70">
                <div className="text-center">
                  <div className="w-12 h-12 mx-auto mb-3 bg-slate-100 group-hover:bg-indigo-100 rounded-xl flex items-center justify-center transition-colors duration-300">
                    <Calendar className="w-6 h-6 text-slate-400 group-hover:text-indigo-600 transition-colors duration-300" />
                  </div>
                  <p className="text-slate-500 font-medium">Aggiungi il tuo evento</p>
                </div>
              </div>
            </div>
          )}
        </div>

        {selectedConference && (
          <div className="fixed bottom-6 right-6 z-50 bg-white rounded-2xl shadow-2xl border border-slate-200 p-6 max-w-sm animate-in slide-in-from-bottom-4">
            <button
              onClick={() => setSelectedConferenceId(null)}
              className="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full bg-slate-100 hover:bg-slate-200 transition-colors"
            >
              <span className="sr-only">Chiudi</span>
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <h3 className="font-bold text-lg text-slate-900 mb-2">{selectedConference.title}</h3>
            <p className="text-sm text-slate-600 mb-4">
              {selectedConference.location} •{" "}
              {new Date(selectedConference.date).toLocaleDateString("it-IT", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })}
            </p>
            <a
              href={selectedConference.website}
              target="_blank"
              rel="noreferrer"
              className="block w-full py-2.5 bg-indigo-600 text-white text-center font-medium rounded-xl hover:bg-indigo-700 transition-colors"
            >
              Vedi dettagli
            </a>
          </div>
        )}

        <section className="mt-20 py-16 px-8 bg-linear-to-br from-slate-900 to-slate-800 rounded-3xl text-center">
          <div className="max-w-2xl mx-auto">
            <Users className="w-12 h-12 text-indigo-400 mx-auto mb-6" />
            <h3 className="text-3xl font-bold text-white mb-4">Unisciti alla community</h3>
            <p className="text-slate-400 text-lg mb-8">
              Migliaia di sviluppatori si sono già uniti per condividere viaggi e esperienze.
            </p>
            {!isAuthenticated ? (
              <Link to="/register" className="inline-block px-8 py-3 bg-white text-slate-900 font-semibold rounded-xl hover:bg-slate-100 transition-all duration-300">
                Registrati gratuitamente
              </Link>
            ) : (
              <p className="text-slate-400">Benvenuto, {user?.name}!</p>
            )}
          </div>
        </section>
      </main>
    </Layout>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/profilo" element={<Profilo />} />
          <Route path="/mie-conferenze" element={<MieConferenze />} />
          <Route path="/conferenze/:conferenceId/registrazione" element={<RegistrazioneConferenza />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
