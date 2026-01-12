import { useState } from "react";
import { useAuth } from "../AuthContext";
import Layout from "../components/Layout";
import { User, MapPin, Edit2, Save, X } from "lucide-react";

export default function Profilo() {
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    name: user?.name || "",
    nickname: user?.nickname || "",
    city: user?.city || "",
    bio: user?.bio || "",
  });

  if (!user) {
    return null;
  }

  const handleSave = () => {
    console.log("Salvando dati:", formData);
    setIsEditing(false);
  };

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-6 py-12">
        <div className="bg-white rounded-3xl shadow-sm border border-slate-200 overflow-hidden">
          <div className="h-32 bg-gradient-to-r from-indigo-500 to-purple-600" />
          
          <div className="px-8 pb-8">
            <div className="relative flex items-end -mt-12 mb-6">
              <div className="w-24 h-24 bg-white rounded-2xl shadow-lg p-1">
                <div className="w-full h-full bg-gradient-to-br from-indigo-100 to-purple-100 rounded-xl flex items-center justify-center">
                  <User className="w-10 h-10 text-indigo-600" />
                </div>
              </div>
              <button
                onClick={() => isEditing ? handleSave() : setIsEditing(true)}
                className="ml-auto px-4 py-2 bg-slate-900 text-white text-sm font-medium rounded-xl hover:bg-slate-800 transition-colors flex items-center gap-2"
              >
                {isEditing ? (
                  <>
                    <Save className="w-4 h-4" />
                    Salva
                  </>
                ) : (
                  <>
                    <Edit2 className="w-4 h-4" />
                    Modifica
                  </>
                )}
              </button>
              {isEditing && (
                <button
                  onClick={() => {
                    setIsEditing(false);
                    setFormData({
                      name: user.name || "",
                      nickname: user.nickname || "",
                      city: user.city || "",
                      bio: user.bio || "",
                    });
                  }}
                  className="ml-2 px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-xl hover:bg-slate-200 transition-colors flex items-center gap-2"
                >
                  <X className="w-4 h-4" />
                  Annulla
                </button>
              )}
            </div>

            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-slate-500 mb-1">Nome</label>
                {isEditing ? (
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="w-full px-4 py-2 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                ) : (
                  <p className="text-lg font-medium text-slate-900">{user.name}</p>
                )}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-slate-500 mb-1">Nickname</label>
                  {isEditing ? (
                    <input
                      type="text"
                      value={formData.nickname}
                      onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
                      className="w-full px-4 py-2 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                  ) : (
                    <p className="text-lg text-slate-900">{user.nickname || "-"}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-500 mb-1">Città</label>
                  {isEditing ? (
                    <input
                      type="text"
                      value={formData.city}
                      onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                      className="w-full px-4 py-2 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                  ) : (
                    <p className="text-lg text-slate-900 flex items-center gap-2">
                      {user.city ? (
                        <>
                          <MapPin className="w-4 h-4 text-slate-400" />
                          {user.city}
                        </>
                      ) : (
                        "-"
                      )}
                    </p>
                  )}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-500 mb-1">Bio</label>
                {isEditing ? (
                  <textarea
                    value={formData.bio}
                    onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                    rows={4}
                    className="w-full px-4 py-2 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 resize-none"
                  />
                ) : (
                  <p className="text-slate-600">{user.bio || "Nessuna bio"}</p>
                )}
              </div>

              <div className="pt-6 border-t border-slate-100">
                <p className="text-sm text-slate-500">Email: {user.email}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}
