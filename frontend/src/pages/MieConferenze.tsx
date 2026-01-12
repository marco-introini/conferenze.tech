import { useState, useEffect } from "react";
import Layout from "../components/Layout";
import ConferenceCard from "../components/ConferenceCard";
import { MapPin, Plus } from "lucide-react";
import type { Conference } from "../types";

export default function MieConferenze() {
  const [conferences, setConferences] = useState<Conference[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadConferences = async () => {
      setIsLoading(false);
      setConferences([]);
    };
    loadConferences();
  }, []);

  return (
    <Layout>
      <div className="max-w-6xl mx-auto px-6 py-12">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Le mie conferenze</h1>
            <p className="text-slate-500 mt-1">Visualizza e gestisci le conferenze a cui partecipi</p>
          </div>
          <button className="px-5 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white text-sm font-semibold rounded-xl hover:shadow-lg hover:shadow-indigo-500/25 transition-all duration-300 flex items-center gap-2">
            <Plus className="w-4 h-4" />
            Aggiungi conferenza
          </button>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center py-20">
            <span className="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin" />
          </div>
        ) : conferences.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {conferences.map((conf) => (
              <ConferenceCard key={conf.id} conf={conf} />
            ))}
          </div>
        ) : (
          <div className="text-center py-20">
            <div className="w-16 h-16 bg-slate-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <MapPin className="w-8 h-8 text-slate-400" />
            </div>
            <h3 className="text-lg font-medium text-slate-900 mb-2">Nessuna conferenza ancora</h3>
            <p className="text-slate-500 mb-6">Non hai ancora aggiunto conferenze al tuo piano</p>
            <button className="px-6 py-3 bg-indigo-600 text-white font-medium rounded-xl hover:bg-indigo-700 transition-colors">
              Esplora conferenze
            </button>
          </div>
        )}
      </div>
    </Layout>
  );
}
