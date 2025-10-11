import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Layout } from './components/Layout'
import { LoginPage } from './pages/LoginPage'
import { DashboardPage } from './pages/DashboardPage'
import { UsersPage } from './pages/UsersPage'
import { AssetsPage } from './pages/AssetsPage'
import { RisksPage } from './pages/RisksPage'
import { DocumentsPage } from './pages/DocumentsPage'
import FileDocumentsPage from './pages/FileDocumentsPage'
import IncidentsPage from './pages/IncidentsPage'
import TrainingPage from './pages/TrainingPage'
import CompliancePage from './pages/CompliancePage'
import AIProvidersPage from './pages/AIProvidersPage'
import AIQueryPage from './pages/AIQueryPage'
import AIChatPage from './pages/AIChatPage'
import RAGManagementPage from './pages/RAGManagementPage'
import RolesManagementPage from './pages/RolesManagementPage'
import OrganizationsPage from './pages/OrganizationsPage'
import { TemplatesPage } from './pages/admin/TemplatesPage'
import { InventoryRulesPage } from './pages/admin/InventoryRulesPage'
import AssetsInventoryPage from './pages/AssetsInventoryPage'

function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/*"
          element={
            <ProtectedRoute>
              <Layout>
                <Routes>
                  <Route path="/" element={<Navigate to="/dashboard" replace />} />
                  <Route path="/dashboard" element={<DashboardPage />} />
                  <Route path="/users" element={<UsersPage />} />
                  <Route path="/admin/roles" element={<RolesManagementPage />} />
                  <Route path="/admin/organizations" element={<OrganizationsPage />} />
                  <Route path="/admin/templates" element={<TemplatesPage />} />
                  <Route path="/admin/inventory-rules" element={<InventoryRulesPage />} />
                  <Route path="/admin/assets-inventory" element={<AssetsInventoryPage />} />
                  <Route path="/assets" element={<AssetsPage />} />
                  <Route path="/risks" element={<RisksPage />} />
                  <Route path="/documents" element={<DocumentsPage />} />
                  <Route path="/file-documents" element={<FileDocumentsPage />} />
                  <Route path="/incidents" element={<IncidentsPage />} />
                  <Route path="/training" element={<TrainingPage />} />
                  <Route path="/compliance" element={<CompliancePage />} />
                  <Route path="/ai/chat" element={<AIChatPage />} />
                  <Route path="/ai/providers" element={<AIProvidersPage />} />
                  <Route path="/ai/query" element={<AIQueryPage />} />
                  <Route path="/ai/rag" element={<RAGManagementPage />} />
                </Routes>
              </Layout>
            </ProtectedRoute>
          }
        />
      </Routes>
    </AuthProvider>
  )
}

export default App
