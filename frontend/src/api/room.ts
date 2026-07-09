import api from "./client";
import type { ApiResponse, Room, RoomSearchParams } from "@/types";

export const roomApi = {
  /** Public: list rooms with filters */
  list(params?: RoomSearchParams) {
    return api.get<ApiResponse<Room[]>>("/rooms", { params });
  },

  /** Public: get room by ID */
  getById(id: string) {
    return api.get<ApiResponse<Room>>(`/rooms/${id}`);
  },

  /** Admin: create room */
  create(data: Partial<Room>) {
    return api.post<ApiResponse<Room>>("/rooms", data);
  },

  /** Admin: update room */
  update(id: string, data: Partial<Room>) {
    return api.put<ApiResponse<Room>>(`/rooms/${id}`, data);
  },

  /** Admin: delete room (soft delete) */
  delete(id: string) {
    return api.delete<ApiResponse>(`/rooms/${id}`);
  },
};
