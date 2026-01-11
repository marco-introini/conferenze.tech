import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import { MapPin, Calendar } from "lucide-react";
import type { Conference } from "../types";
import "leaflet/dist/leaflet.css";
import L from "leaflet";

interface Props {
  conferences: Conference[];
  onSelectConference: (id: string) => void;
}

const ITALY_CENTER: [number, number] = [42.5, 12.5];
const DEFAULT_ZOOM = 6;

const icon = L.icon({
  iconUrl: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%236366f1' stroke='white' stroke-width='2'%3E%3Cpath d='M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z'/%3E%3Ccircle cx='12' cy='10' r='3' fill='white'/%3E%3C/svg%3E",
  iconSize: [32, 32],
  iconAnchor: [16, 32],
  popupAnchor: [0, -32],
});

const cityCoords: Record<string, [number, number]> = {
  "Firenze (FI)": [43.7696, 11.2558],
  "Verona (VR)": [45.4384, 10.9916],
  "Milano (MI)": [45.4642, 9.19],
  "Roma (RM)": [41.9028, 12.4964],
  "Bologna (BO)": [44.4949, 11.3426],
  "Torino (TO)": [45.0703, 7.6869],
  "Napoli (NA)": [40.8518, 14.2681],
  "Genova (GE)": [44.4056, 8.9463],
  "Palermo (PA)": [38.1157, 13.3615],
  "Catania (CT)": [37.5079, 15.0833],
};

function getCoords(conf: Conference): [number, number] {
  if (conf.latitude && conf.longitude) {
    return [conf.latitude, conf.longitude];
  }
  return cityCoords[conf.location] || ITALY_CENTER;
}

export default function ConferenceMap({ conferences, onSelectConference }: Props) {
  const validConferences = conferences.filter(c => {
    if (c.latitude && c.longitude) return true;
    return c.location in cityCoords;
  });

  return (
    <div className="relative w-full h-[500px] rounded-2xl overflow-hidden border border-slate-200 shadow-lg">
      <MapContainer
        center={ITALY_CENTER}
        zoom={DEFAULT_ZOOM}
        className="w-full h-full"
        scrollWheelZoom={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {validConferences.map((conf) => {
          const coords = getCoords(conf);
          return (
            <Marker
              key={conf.id}
              position={coords}
              icon={icon}
              eventHandlers={{
                click: () => onSelectConference(conf.id),
              }}
            >
              <Popup className="custom-popup">
                <div className="p-2 min-w-[200px]">
                  <h3 className="font-bold text-slate-800 mb-2">{conf.title}</h3>
                  <div className="flex items-center gap-2 text-sm text-slate-600 mb-1">
                    <Calendar className="w-4 h-4" />
                    <span>
                      {new Date(conf.date).toLocaleDateString("it-IT", {
                        day: "numeric",
                        month: "short",
                      })}
                    </span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-slate-600">
                    <MapPin className="w-4 h-4" />
                    <span>{conf.location}</span>
                  </div>
                  <button
                    onClick={() => onSelectConference(conf.id)}
                    className="mt-3 w-full px-3 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors"
                  >
                    Vedi dettagli
                  </button>
                </div>
              </Popup>
            </Marker>
          );
        })}
      </MapContainer>
    </div>
  );
}
