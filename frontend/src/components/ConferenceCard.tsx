import {
  Calendar,
  MapPin,
  ExternalLink,
  Car,
  User as UserIcon,
} from "lucide-react";
import { useNavigate } from "react-router-dom";
import type { Conference } from "../types";

interface Props {
  conf: Conference;
}

export default function ConferenceCard({ conf }: Props) {
  const navigate = useNavigate();
  const attendees = conf.attendees || [];

  return (
    <div className="group bg-white rounded-2xl border border-slate-200 overflow-hidden hover:shadow-xl hover:shadow-slate-200/50 hover:-translate-y-1 transition-all duration-300">
      <div className="p-6">
        <div className="flex justify-between items-start mb-4">
          <h3 className="text-xl font-bold text-slate-800 group-hover:text-indigo-600 transition-colors">{conf.title}</h3>
          <span className="bg-indigo-50 text-indigo-700 text-xs font-semibold px-3 py-1 rounded-full">
            Tech
          </span>
        </div>

        <div className="space-y-3 text-sm">
          <div className="flex items-center gap-3 text-slate-600">
            <div className="w-8 h-8 bg-slate-100 rounded-lg flex items-center justify-center">
              <Calendar className="w-4 h-4 text-slate-500" />
            </div>
            <span>
              {new Date(conf.date).toLocaleDateString("it-IT", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })}
            </span>
          </div>
          <div className="flex items-center gap-3 text-slate-600">
            <div className="w-8 h-8 bg-slate-100 rounded-lg flex items-center justify-center">
              <MapPin className="w-4 h-4 text-slate-500" />
            </div>
            <span>{conf.location}</span>
          </div>
          {conf.website && (
            <a
              href={conf.website}
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-2 text-indigo-600 hover:text-indigo-700 transition-colors text-sm font-medium"
            >
              <ExternalLink className="w-4 h-4" />
              <span>Sito Ufficiale</span>
            </a>
          )}
        </div>
      </div>

        <div className="px-6 py-4 bg-slate-50 border-t border-slate-100 flex items-center justify-between">
          <span className="text-sm font-medium text-slate-500">
            {attendees.length} partecipanti
          </span>
          <button
            onClick={() => navigate(`/conferenze/${conf.id}/registrazione`)}
            className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors"
          >
            Partecipo
          </button>
        </div>

      <div className="p-6 pt-4">
        <h4 className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-4">
          Community & Passaggi
        </h4>
        {attendees.length > 0 ? (
          <div className="space-y-3">
            {attendees.map((att, idx) => (
              <div
                key={idx}
                className="flex items-center justify-between text-sm p-3 rounded-xl hover:bg-slate-50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="w-9 h-9 rounded-full bg-gradient-to-br from-slate-100 to-slate-200 flex items-center justify-center text-slate-500">
                    <UserIcon className="w-4 h-4" />
                  </div>
                  <div>
                    <p className="font-medium text-slate-700">
                      {att.user.nickname || "Utente"}
                    </p>
                    {att.user.city && (
                      <p className="text-xs text-slate-500">da {att.user.city}</p>
                    )}
                  </div>
                </div>

                <div className="flex gap-1.5">
                  {att.needsRide && (
                    <span
                      title="Cerca passaggio"
                      className="p-2 bg-amber-100 text-amber-700 rounded-lg cursor-help"
                    >
                      <Car className="w-3.5 h-3.5" />
                    </span>
                  )}
                  {att.hasCar && (
                    <span
                      title="Offre passaggio"
                      className="p-2 bg-emerald-100 text-emerald-700 rounded-lg cursor-help"
                    >
                      <Car className="w-3.5 h-3.5" />
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-slate-500 text-center py-4">
            Nessun partecipante ancora
          </p>
        )}
      </div>
    </div>
  );
}
