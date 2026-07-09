import type { Hotel } from "./hotel";

export interface RoomImage {
  roomId: string;
  imageUrl: string;
}

export interface Room {
  id: string;
  hotelId: string;
  typeName: string;
  price: number;
  totalQuantity: number;
  availableQuantity: number;
  description?: string;
  images: RoomImage[];
  createAt: string;
  updateAt: string;
  status: number;
  hotel?: Hotel;
}

export interface RoomDetails {
  roomId: string;
  hotelId: string;
  hotelName: string;
  typeName: string;
  price: number;
  totalQuantity: number;
  availableQuantity: number;
  roomDescription?: string;
  province?: string;
  city?: string;
  district?: string;
}

export interface RoomSearchParams {
  page?: number;
  pageSize?: number;
  hotelID?: string;
}
