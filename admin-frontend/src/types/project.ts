export interface Project {
  id: number;
  name: string;
  description: string;
  slug: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface ProjectFormData {
  name: string;
  description: string;
  slug: string;
}