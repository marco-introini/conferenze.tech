import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Layout from "../components/Layout";
import { api, type Conference } from "../api";
import { useAuth } from "../AuthContext";
import { Car, User, ArrowLeft, Check } from "lucide-react";

export default function RegistrazioneConferenza() {
  const navigate = useNavigate();
  const { conferenceId } = useParams();
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const [conference, setConference] = useState<Conference | null>(null);
  const [needsRide, setNeedsRide] = useState(false);
  const [hasCar, setHasCar] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!isAuthLoading && !isAuthenticated) {
      navigate("/login");
    }
  }, [isAuthLoading, isAuthenticated, navigate]);

  useEffect(() => {
    const loadConference = async () => {
      if (!conferenceId) {
        console.log("[RegistrazioneConferenza] conferenceId mancante");
        setIsLoading(false);
        return;
      }
      console.log("[RegistrazioneConferenza] Caricamento conferenza:", conferenceId);
      try {
        const data = await api.getConference(conferenceId);
        console.log("[RegistrazioneConferenza] Conferenza caricata:", data);
        setConference(data);
      } catch (err) {
        console.error("[RegistrazioneConferenza] Errore nel caricamento conferenza:", err);
        setError("Impossibile caricare i dettagli della conferenza");
      } finally {
        setIsLoading(false);
      }
    };

    loadConference();
  }, [conferenceId]);

  const handleSubmit = async () => {
    if (!conference) return;
    setIsSubmitting(true);
    setError("");

    try {
      await api.registerToConference(conference.id, {
        role: "attendee",
        needsRide,
        hasCar,
      });
      navigate("/mie-conferenze");
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Errore nella registrazione";
      if (msg.includes("già registrato") || msg.includes("already registered")) {
        setError("Sei già registrato a questa conferenza");
      } else {
        setError(msg);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isAuthLoading || isLoading) {
    return (
      <Layout>
        <div className="flex items-center justify-center py-32">
          <span className="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin" />
        </div>
      </Layout>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  if (!conference) {
    return (
      <Layout>
        <div className="max-w-md mx-auto px-6 py-12 text-center">
          <p className="text-slate-600 mb-6">{error || "Conferenza non trovata"}</p>
          <button
            onClick={() => navigate("/")}
            className="px-6 py-3 bg-indigo-600 text-white font-medium rounded-xl hover:bg-indigo-700 transition-colors"
          >
            Torna alla home
          </button>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <main className="pt-32 pb-20 px-6">
        <div className="max-w-lg mx-auto">
          <button
            onClick={() => navigate(-1)}
            className="flex items-center gap-2 text-slate-600 hover:text-slate-900 mb-8 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
            <span className="text-sm font-medium">Indietro</span>
          </button>

          <div className="bg-white rounded-2xl shadow-xl border border-slate-200 p-8">
            <div className="mb-8">
              <h1 className="text-2xl font-bold text-slate-900 mb-2">Registrati alla conferenza</h1>
              <p className="text-slate-600">Compila i dettagli per confermare la tua partecipazione</p>
            </div>

            <div className="bg-slate-50 rounded-xl p-5 mb-8">
              <h2 className="font-semibold text-slate-900 mb-1">{conference.title}</h2>
              <p className="text-sm text-slate-600">
                {new Date(conference.date).toLocaleDateString("it-IT", {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                })}{" "}
                • {conference.location}
              </p>
            </div>

            {error && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
                {error}
              </div>
            )}

            <div className="space-y-5 mb-8">
              <div className="flex items-start gap-4 p-4 bg-white border border-slate-200 rounded-xl">
                <div className="w-10 h-10 bg-amber-100 text-amber-700 rounded-lg flex items-center justify-center flex-shrink-0">
                  <Car className="w-5 h-5" />
                </div>
                <div className="flex-1">
                  <label className="block font-medium text-slate-900 mb-1">Hai bisogno di un passaggio?</label>
                  <p className="text-sm text-slate-600 mb-3">
                    Se non hai un mezzo per raggiungere la conferenza, indica che cerchi un passaggio
                  </p>
                  <div className="flex gap-3">
                    <button
                      type="button"
                      onClick={() => setNeedsRide(true)}
                      className={`flex-1 py-2.5 rounded-lg text-sm font-medium transition-all ${
                        needsRide
                          ? "bg-amber-100 text-amber-800 border-2 border-amber-500"
                          : "bg-slate-100 text-slate-600 border-2 border-transparent hover:bg-slate-200"
                      }`}
                    >
                      Sì, cerco passaggio
                    </button>
                    <button
                      type="button"
                      onClick={() => setNeedsRide(false)}
                      className={`flex-1 py-2.5 rounded-lg text-sm font-medium transition-all ${
                        !needsRide
                          ? "bg-slate-100 text-slate-800 border-2 border-slate-400"
                          : "bg-white text-slate-600 border-2 border-slate-200 hover:bg-slate-50"
                      }`}
                    >
                      No, ce l'ho
                    </button>
                  </div>
                </div>
              </div>

              <div className="flex items-start gap-4 p-4 bg-white border border-slate-200 rounded-xl">
                <div className="w-10 h-10 bg-emerald-100 text-emerald-700 rounded-lg flex items-center justify-center flex-shrink-0">
                  <User className="w-5 h-5" />
                </div>
                <div className="flex-1">
                  <label className="block font-medium text-slate-900 mb-1">Hai una macchina disponibile?</label>
                  <p className="text-sm text-slate-600 mb-3">
                    Se hai un'auto e vuoi offrire passaggi ad altri partecipanti, seleziona Sì
                  </p>
                  <div className="flex gap-3">
                    <button
                      type="button"
                      onClick={() => setHasCar(true)}
                      className={`flex-1 py-2.5 rounded-lg text-sm font-medium transition-all ${
                        hasCar
                          ? "bg-emerald-100 text-emerald-800 border-2 border-emerald-500"
                          : "bg-slate-100 text-slate-600 border-2 border-transparent hover:bg-slate-200"
                      }`}
                    >
                      Sì, offro passaggi
                    </button>
                    <button
                      type="button"
                      onClick={() => setHasCar(false)}
                      className={`flex-1 py-2.5 rounded-lg text-sm font-medium transition-all ${
                        !hasCar
                          ? "bg-slate-100 text-slate-800 border-2 border-slate-400"
                          : "bg-white text-slate-600 border-2 border-slate-200 hover:bg-slate-50"
                      }`}
                    >
                      No
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <button
              onClick={handleSubmit}
              disabled={isSubmitting}
              className="w-full py-3.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:shadow-lg hover:shadow-indigo-500/25 transition-all duration-300 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? (
                <span className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              ) : (
                <>
                  <Check className="w-4 h-4" />
                  Conferma registrazione
                </>
              )}
            </button>
          </div>
        </div>
      </main>
    </Layout>
  );
}
