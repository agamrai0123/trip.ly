import type { City, TripPost } from "@/types";
import tokyo from "@/assets/city-tokyo.jpg";
import santorini from "@/assets/city-santorini.jpg";
import bali from "@/assets/city-bali.jpg";
import lisbon from "@/assets/city-lisbon.jpg";
import kyoto from "@/assets/city-kyoto.jpg";
import reykjavik from "@/assets/city-reykjavik.jpg";

export const cities: City[] = [
  { id: "tokyo", name: "Tokyo", country: "Japan", image: tokyo, blurb: "Neon nights, tranquil shrines.", postsCount: 184 },
  { id: "santorini", name: "Santorini", country: "Greece", image: santorini, blurb: "Whitewashed cliffs over the Aegean.", postsCount: 122 },
  { id: "bali", name: "Bali", country: "Indonesia", image: bali, blurb: "Rice terraces and ocean temples.", postsCount: 209 },
  { id: "lisbon", name: "Lisbon", country: "Portugal", image: lisbon, blurb: "Trams, tiles, and Atlantic light.", postsCount: 96 },
  { id: "kyoto", name: "Kyoto", country: "Japan", image: kyoto, blurb: "Bamboo groves and quiet tea houses.", postsCount: 158 },
  { id: "reykjavik", name: "Reykjavik", country: "Iceland", image: reykjavik, blurb: "Aurora, geysers, the wild north.", postsCount: 71 },
];

const u = (seed: string) => `https://api.dicebear.com/9.x/notionists/svg?seed=${seed}`;
const ph = (q: string, w = 800, h = 600) => `https://images.unsplash.com/photo-${q}?w=${w}&h=${h}&fit=crop&auto=format&q=80`;

export const tripPosts: TripPost[] = [
  {
    id: "p1",
    cityId: "tokyo",
    author: { name: "Maya Chen", avatar: u("Maya") },
    title: "Tokyo in 5 days — neon, ramen, and quiet shrines",
    days: 5,
    totalBudget: 1840,
    bestTime: "March – May (cherry blossoms)",
    likes: 482,
    cover: ph("1542051841857-5f90071e7989", 1200, 700),
    summary: "A balanced loop between Shibuya energy and Yanaka calm. Walkable, transit-friendly, with two splurge meals.",
    places: [
      { id: "p1-1", name: "Sensō-ji Temple", category: "sight", slot: "morning", lat: 35.7148, lng: 139.7967, address: "Asakusa, Taitō", expense: 0, review: "Arrive before 8am — the lantern is yours alone.", rating: 4.8, bestTime: "Sunrise", image: ph("1492571350019-22de08371fd3") },
      { id: "p1-2", name: "Afuri Ramen Ebisu", category: "food", slot: "lunch", lat: 35.6464, lng: 139.7100, address: "Ebisu", expense: 14, review: "Yuzu shio broth — light, citrusy, crave-worthy.", rating: 4.7, image: ph("1569718212165-3a8278d5f624") },
      { id: "p1-3", name: "teamLab Planets", category: "activity", slot: "evening", lat: 35.6493, lng: 139.7903, address: "Toyosu", expense: 32, review: "Wear shorts. The mirror room is unreal.", rating: 4.6, image: ph("1493514789931-586cb221d7a7") },
      { id: "p1-4", name: "Omoide Yokocho", category: "food", slot: "night", lat: 35.6938, lng: 139.6996, address: "Shinjuku", expense: 28, review: "Tiny yakitori alleys. Skewers + highballs.", rating: 4.5, image: ph("1554188248-986adbb73be4") },
      { id: "p1-5", name: "The Knot Shinjuku", category: "stay", slot: "night", lat: 35.6938, lng: 139.6900, address: "Nishi-Shinjuku", expense: 142, review: "Compact rooms, killer rooftop views.", rating: 4.4, image: ph("1566073771259-6a8506099945") },
    ],
    comments: [
      { id: "c1", author: "Diego", avatar: u("Diego"), text: "Saved! Adding teamLab to my September trip.", at: "2d" },
      { id: "c2", author: "Priya", avatar: u("Priya"), text: "Afuri yuzu shio is unreal. Try the tsukemen too.", at: "5d" },
    ],
  },
  {
    id: "p2",
    cityId: "tokyo",
    author: { name: "Leo Park", avatar: u("Leo") },
    title: "Budget Tokyo: 4 days under $900",
    days: 4,
    totalBudget: 870,
    bestTime: "October (mild, low rain)",
    likes: 311,
    cover: ph("1503899036084-c55cdd92da26", 1200, 700),
    summary: "Hostels, conbini breakfasts, and free parks. Still hit every must-see.",
    places: [
      { id: "p2-1", name: "Meiji Jingu", category: "sight", slot: "morning", lat: 35.6764, lng: 139.6993, address: "Shibuya", expense: 0, review: "Forest in the city. Free and humbling.", rating: 4.9, image: ph("1480796927426-f609979314bd") },
      { id: "p2-2", name: "Ichiran Shibuya", category: "food", slot: "lunch", lat: 35.6595, lng: 139.7005, address: "Shibuya", expense: 12, review: "Solo booth ramen. Touristy but legit.", rating: 4.3, image: ph("1591814468924-caf88d1232e1") },
      { id: "p2-3", name: "Shibuya Sky", category: "sight", slot: "evening", lat: 35.6587, lng: 139.7016, address: "Shibuya Scramble Sq.", expense: 18, review: "Book sunset slot 3 weeks ahead.", rating: 4.8, image: ph("1540959733332-eab4deabeeaf") },
      { id: "p2-4", name: "UNPLAN Kagurazaka", category: "stay", slot: "night", lat: 35.7022, lng: 139.7384, address: "Kagurazaka", expense: 38, review: "Best hostel in the city. Bunk pods.", rating: 4.7, image: ph("1551776245-d7a86d8cb9e6") },
    ],
    comments: [
      { id: "c3", author: "Sana", avatar: u("Sana"), text: "Doing this verbatim next month 🙏", at: "1w" },
    ],
  },
  {
    id: "p3",
    cityId: "santorini",
    author: { name: "Élise Moreau", avatar: u("Elise") },
    title: "Santorini slow: 6 days of caldera + cave wines",
    days: 6,
    totalBudget: 2450,
    bestTime: "May or September (no crowds)",
    likes: 528,
    cover: ph("1570077188670-e3a8d69ac5ff", 1200, 700),
    summary: "Skip Oia at sunset. Eat where the locals eat. Rent an ATV.",
    places: [
      { id: "p3-1", name: "Akrotiri ruins", category: "sight", slot: "morning", lat: 36.3517, lng: 25.4036, address: "Akrotiri", expense: 12, review: "Bronze Age Pompeii. Empty by 9am.", rating: 4.7, bestTime: "Spring", image: ph("1533104816931-20fa691ff6ca") },
      { id: "p3-2", name: "To Psaraki", category: "food", slot: "lunch", lat: 36.3499, lng: 25.4419, address: "Vlychada Marina", expense: 42, review: "Fresh-off-the-boat seafood. Order the calamari.", rating: 4.8, image: ph("1559339352-11d035aa65de") },
      { id: "p3-3", name: "Santo Wines tasting", category: "activity", slot: "evening", lat: 36.3909, lng: 25.4625, address: "Pyrgos", expense: 36, review: "Caldera-edge sunset with assyrtiko flight.", rating: 4.6, image: ph("1474722883778-792e7990302f") },
      { id: "p3-4", name: "Cave House Imerovigli", category: "stay", slot: "night", lat: 36.4314, lng: 25.4322, address: "Imerovigli", expense: 280, review: "Plunge pool over the caldera. Worth it.", rating: 4.9, image: ph("1582719508461-905c673771fd") },
    ],
    comments: [
      { id: "c4", author: "Theo", avatar: u("Theo"), text: "Akrotiri >> Oia for sunset. This.", at: "3d" },
    ],
  },
  {
    id: "p4",
    cityId: "santorini",
    author: { name: "Noor Aziz", avatar: u("Noor") },
    title: "3 days in Oia for first-timers",
    days: 3,
    totalBudget: 1180,
    bestTime: "Late May",
    likes: 196,
    cover: ph("1613395877344-13d4a8e0d49e", 1200, 700),
    summary: "Tight 3-day loop hitting the icons without burning out.",
    places: [
      { id: "p4-1", name: "Oia Castle viewpoint", category: "sight", slot: "evening", lat: 36.4618, lng: 25.3753, address: "Oia", expense: 0, review: "Get there 90 min before sunset. Bring water.", rating: 4.5, image: ph("1469474968028-56623f02e42e") },
      { id: "p4-2", name: "Lolita's Gelato", category: "food", slot: "lunch", lat: 36.4612, lng: 25.3756, address: "Oia main path", expense: 6, review: "Pistachio + fig. Twice.", rating: 4.9, image: ph("1567206563064-6f60f40a2b57") },
      { id: "p4-3", name: "Canaves Oia Suites", category: "stay", slot: "night", lat: 36.4625, lng: 25.3758, address: "Oia", expense: 410, review: "Splurge stay — infinity pool, white-on-white.", rating: 4.8, image: ph("1551918120-9739cb430c6d") },
    ],
    comments: [],
  },
  {
    id: "p5",
    cityId: "bali",
    author: { name: "Arjun Rao", avatar: u("Arjun") },
    title: "Bali surf + jungle: 8 days Canggu to Ubud",
    days: 8,
    totalBudget: 1320,
    bestTime: "April – October (dry season)",
    likes: 612,
    cover: ph("1537996194471-e657df975ab4", 1200, 700),
    summary: "Two bases, zero scooter accidents. Surf mornings, jungle afternoons.",
    places: [
      { id: "p5-1", name: "Batu Bolong beach", category: "activity", slot: "morning", lat: -8.6553, lng: 115.1311, address: "Canggu", expense: 8, review: "Rent a foamie for $5. Mellow lefts.", rating: 4.6, image: ph("1507525428034-b723cf961d3e") },
      { id: "p5-2", name: "Crate Café", category: "food", slot: "lunch", lat: -8.6542, lng: 115.1366, address: "Canggu", expense: 9, review: "Acai + smashed avo. Worth the queue.", rating: 4.5, image: ph("1517248135467-4c7edcad34c4") },
      { id: "p5-3", name: "Tegallalang rice terraces", category: "sight", slot: "morning", lat: -8.4317, lng: 115.2785, address: "Ubud", expense: 4, review: "Go before 8 to beat tour buses.", rating: 4.7, image: ph("1518391846015-55a9cc003b25") },
      { id: "p5-4", name: "Bambu Indah", category: "stay", slot: "night", lat: -8.4985, lng: 115.2625, address: "Ubud", expense: 220, review: "Bamboo eco-villas in a river valley.", rating: 4.9, image: ph("1540541338287-41700207dee6") },
    ],
    comments: [
      { id: "c5", author: "Mei", avatar: u("Mei"), text: "Crate is unreal. Add Penny Lane too 🍳", at: "4d" },
    ],
  },
  {
    id: "p6",
    cityId: "lisbon",
    author: { name: "Joana Silva", avatar: u("Joana") },
    title: "Lisbon weekend: pastéis, fado, miradouros",
    days: 3,
    totalBudget: 540,
    bestTime: "April or October",
    likes: 274,
    cover: ph("1513735492246-483525079686", 1200, 700),
    summary: "Hilly walks, tile-everything, sardines under fado guitars.",
    places: [
      { id: "p6-1", name: "Miradouro da Senhora do Monte", category: "sight", slot: "morning", lat: 38.7188, lng: -9.1325, address: "Graça", expense: 0, review: "Best free view in the city. Quiet at 8am.", rating: 4.8, image: ph("1518733057094-95b53143d2a7") },
      { id: "p6-2", name: "Time Out Market", category: "food", slot: "lunch", lat: 38.7071, lng: -9.1456, address: "Cais do Sodré", expense: 22, review: "One stop, every great chef. Gets packed.", rating: 4.5, image: ph("1555396273-367ea4eb4db5") },
      { id: "p6-3", name: "Fado in Alfama", category: "activity", slot: "night", lat: 38.7117, lng: -9.1304, address: "Alfama", expense: 45, review: "Tiny tasca, two singers. Goosebumps.", rating: 4.9, image: ph("1485872299712-a8e4f4b3a6c2") },
      { id: "p6-4", name: "Memmo Alfama", category: "stay", slot: "night", lat: 38.7115, lng: -9.1308, address: "Alfama", expense: 165, review: "Rooftop pool over terracotta roofs.", rating: 4.7, image: ph("1566073771259-6a8506099945") },
    ],
    comments: [],
  },
  {
    id: "p7",
    cityId: "kyoto",
    author: { name: "Hiro Tanaka", avatar: u("Hiro") },
    title: "Kyoto in autumn: 4 days of momiji",
    days: 4,
    totalBudget: 1050,
    bestTime: "Mid-November",
    likes: 389,
    cover: ph("1528360983277-13d401cdc186", 1200, 700),
    summary: "Temple loops timed for peak red leaves. One ryokan night.",
    places: [
      { id: "p7-1", name: "Fushimi Inari", category: "sight", slot: "morning", lat: 34.9671, lng: 135.7727, address: "Fushimi", expense: 0, review: "Hike past the crowds — top is empty.", rating: 4.9, bestTime: "Sunrise", image: ph("1528360983277-13d401cdc186") },
      { id: "p7-2", name: "Nishiki Market", category: "food", slot: "lunch", lat: 35.0050, lng: 135.7641, address: "Nakagyō", expense: 18, review: "Sample-as-you-go. Try tako tamago.", rating: 4.6, image: ph("1591814468924-caf88d1232e1") },
      { id: "p7-3", name: "Tofukuji at sunset", category: "sight", slot: "evening", lat: 34.9764, lng: 135.7740, address: "Higashiyama", expense: 8, review: "Maple bridge. November is unreal.", rating: 4.8, image: ph("1545569341-9eb8b30979d9") },
      { id: "p7-4", name: "Tawaraya Ryokan", category: "stay", slot: "night", lat: 35.0116, lng: 135.7681, address: "Nakagyō", expense: 480, review: "300-year-old ryokan. Kaiseki dinner included.", rating: 5.0, image: ph("1542640244-7e672d6cef4e") },
    ],
    comments: [
      { id: "c6", author: "Aiko", avatar: u("Aiko"), text: "Tawaraya was a once-in-a-lifetime night. 🥹", at: "1d" },
    ],
  },
  {
    id: "p8",
    cityId: "reykjavik",
    author: { name: "Sigrid Holm", avatar: u("Sigrid") },
    title: "Reykjavik + Golden Circle in 4 days",
    days: 4,
    totalBudget: 1680,
    bestTime: "September – March (auroras)",
    likes: 217,
    cover: ph("1490093158370-1d3f4cc1c2a6", 1200, 700),
    summary: "City base + day trips. Aurora hunt on night 2.",
    places: [
      { id: "p8-1", name: "Hallgrímskirkja", category: "sight", slot: "morning", lat: 64.1418, lng: -21.9266, address: "Reykjavik", expense: 8, review: "Tower view of the colored rooftops.", rating: 4.6, image: ph("1531168556467-80aace0d0144") },
      { id: "p8-2", name: "Sægreifinn (Sea Baron)", category: "food", slot: "lunch", lat: 64.1521, lng: -21.9432, address: "Old Harbour", expense: 24, review: "Lobster soup. Shack vibes, Michelin flavor.", rating: 4.7, image: ph("1565299624946-b28f40a0ae38") },
      { id: "p8-3", name: "Þingvellir aurora hunt", category: "activity", slot: "night", lat: 64.2559, lng: -21.1295, address: "Þingvellir NP", expense: 0, review: "Drive 45 min out. Check kp index. Patience.", rating: 4.8, bestTime: "Oct–Mar", image: ph("1483347756197-71ef80e95f73") },
      { id: "p8-4", name: "ION City Hotel", category: "stay", slot: "night", lat: 64.1480, lng: -21.9408, address: "Laugavegur", expense: 195, review: "Design-y, central, warm.", rating: 4.5, image: ph("1582719508461-905c673771fd") },
    ],
    comments: [],
  },
];
