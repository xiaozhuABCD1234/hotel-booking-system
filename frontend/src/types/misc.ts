export interface Region {
  id: number;
  regionName: string;
  parentsId?: number;
}

export interface GuestStats {
  ageGroup: string;
  gender: string;
  totalOrders: number;
  totalAmount: number;
}

export interface PersonInfo {
  idCard: string;
  name: string;
  gender: string;
  age: number;
}

export interface GuestBookingStats {
  personIdCard: string;
  personName: string;
  gender: string;
  age: number;
  totalOrders: number;
  totalAmount: number;
}
