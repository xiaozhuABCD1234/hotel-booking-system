import type { Region } from "./misc";

export interface Hotel {
  id: string;
  hotelName: string;
  regionId: number;
  region?: Region;
  address: string;
  telephone: string;
  starLevel?: number;
  description?: string;
  createAt: string;
  updateAt: string;
  status: number;
  images: HotelImage[];
}

export interface HotelImage {
  hotelId: string;
  imageUrl: string;
}

export interface HotelSummary {
  hotelId: string;
  hotelName: string;
  province?: string;
  city?: string;
  district: string;
  address: string;
  telephone: string;
  starLevel?: number;
  description?: string;
  mainImage?: string;
  minPrice: number;
  roomCount: number;
  totalRooms: number;
  avgRating: number;
  reviewCount: number;
  status: number;
}

export interface HotelSearchParams {
  page?: number;
  pageSize?: number;
  regionID?: number;
  starLevel?: number;
  keyword?: string;
  minPrice?: number;
  maxPrice?: number;
  checkInDate?: string;
  checkOutDate?: string;
}
