export interface Region {
  id: number
  name: string
  parentId?: number
  level: number
}

export interface GuestStats {
  ageGroup: string
  gender: string
  occupation?: string
  education?: string
  income?: string
  bookingCount: number
  totalSpent: number
}

export interface PersonInfo {
  idCard: string
  name: string
  gender: string
  age: number
  occupation?: string
  education?: string
  income?: string
}

export interface GuestBookingStats {
  personId: string
  name: string
  gender: string
  age: number
  occupation?: string
  bookingCount: number
  totalSpent: number
}
