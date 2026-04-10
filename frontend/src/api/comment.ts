import axios from 'axios'
import type { ApiResponse } from './project'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 发送cookies
})

// 评论用户信息
export interface CommentUser {
  id: number
  name: string
  avatar: string
}

// 评论响应
export interface Comment {
  id: number
  content: string
  user: CommentUser
  task_id: number
  created_at: string
  updated_at: string
}

// 创建评论请求
export interface CreateCommentRequest {
  content: string
}

// 更新评论请求
export interface UpdateCommentRequest {
  content: string
}

// 获取任务的评论列表
export function getComments(taskId: number): Promise<ApiResponse<Comment[]>> {
  return apiClient.get(`/tasks/${taskId}/comments`).then(res => res.data)
}

// 添加评论
export function createComment(taskId: number, content: string): Promise<ApiResponse<Comment>> {
  return apiClient.post(`/tasks/${taskId}/comments`, { content }).then(res => res.data)
}

// 更新评论
export function updateComment(taskId: number, commentId: number, content: string): Promise<ApiResponse<Comment>> {
  return apiClient.put(`/tasks/${taskId}/comments/${commentId}`, { content }).then(res => res.data)
}

// 删除评论
export function deleteComment(taskId: number, commentId: number): Promise<ApiResponse<null>> {
  return apiClient.delete(`/tasks/${taskId}/comments/${commentId}`).then(res => res.data)
}
