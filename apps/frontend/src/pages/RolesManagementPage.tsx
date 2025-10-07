import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Box,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Tooltip,
  CircularProgress,
  Chip,
  Checkbox,
  FormControlLabel,
  FormGroup,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import {
  Add,
  Edit,
  Delete,
  Security,
  ExpandMore,
} from '@mui/icons-material';
import { rolesApi, Role, RoleWithPermissions, Permission } from '../shared/api/roles';
import { useAuth } from '../contexts/AuthContext';

// Функция для получения иконки модуля
const getModuleIcon = (module: string) => {
  const icons: Record<string, string> = {
    'Документы': '📄',
    'Инциденты': '⚠️',
    'Риски': '📊',
    'Активы': '🛠️',
    'Обучение': '🎓',
    'ИИ': '🤖',
    'Пользователи': '👤',
    'Роли': '🛡️',
    'Аудит': '📜',
    'Дашборд': '📈',
    'Соответствие': '✅',
    'Отчеты': '📋'
  };
  return icons[module] || '📁';
};

const RolesManagementPage: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [roleUserCounts, setRoleUserCounts] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(true);
  const [selectedRole, setSelectedRole] = useState<RoleWithPermissions | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showPermissionsModal, setShowPermissionsModal] = useState(false);

  const { user } = useAuth();

  useEffect(() => {
    if (user) {
      loadData();
    } else {
      setLoading(false);
    }
  }, [user]);

  const loadData = async () => {
    try {
      setLoading(true);
      const [rolesRaw, permissionsRaw] = await Promise.all([
        rolesApi.getRoles(),
        rolesApi.getPermissions()
      ]);
      // Нормализация ключей в camelCase на случай, если пришли поля в PascalCase
      const rolesData: Role[] = (rolesRaw as any[]).map((r) => ({
        id: r.id ?? r.ID,
        name: r.name ?? r.Name,
        description: r.description ?? r.Description ?? undefined,
        created_at: r.created_at ?? r.CreatedAt,
        updated_at: r.updated_at ?? r.UpdatedAt,
      }));
      const permissionsData: Permission[] = (permissionsRaw as any[]).map((p) => ({
        id: p.id ?? p.ID,
        code: p.code ?? p.Code,
        module: p.module ?? p.Module,
        description: p.description ?? p.Description ?? undefined,
      }));
      setRoles(rolesData);
      setPermissions(permissionsData);
      
      // Загружаем количество пользователей для каждой роли
      const userCounts: Record<string, number> = {};
      for (const role of rolesData) {
        try {
          const users = await rolesApi.getRoleUsers(role.id);
          userCounts[role.id] = users.length;
        } catch (error) {
          console.error(`Ошибка загрузки пользователей для роли ${role.name}:`, error);
          userCounts[role.id] = 0;
        }
      }
      setRoleUserCounts(userCounts);
      
      // Debug: проверим структуру данных
      console.log('Загруженные роли (нормализовано):', rolesData);
      console.log('Пример роли:', rolesData[0]);
      console.log('Загруженные права (нормализовано):', permissionsData);
      console.log('Пример права:', permissionsData[0]);
      console.log('Количество пользователей по ролям:', userCounts);
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateRole = async (roleData: any) => {
    try {
      await rolesApi.createRole(roleData);
      await loadData();
      setShowCreateModal(false);
    } catch (error) {
      console.error('Ошибка создания роли:', error);
    }
  };

  const handleUpdateRole = async (id: string, roleData: any) => {
    try {
      console.log('handleUpdateRole called with:', { id, roleData });
      console.log('roleData.permission_ids:', roleData.permission_ids);
      console.log('roleData.permission_ids length:', roleData.permission_ids?.length);
      
      await rolesApi.updateRole(id, roleData);
      console.log('Role updated successfully');
      await loadData();
      setShowEditModal(false);
      setSelectedRole(null);
    } catch (error) {
      console.error('Ошибка обновления роли:', error);
    }
  };

  const handleDeleteRole = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить роль?')) {
      try {
        await rolesApi.deleteRole(id);
        await loadData();
      } catch (error) {
        console.error('Ошибка удаления роли:', error);
      }
    }
  };

  const handleViewPermissions = async (role: Role) => {
    try {
      const roleWithPermissions = await rolesApi.getRole(role.id);
      setSelectedRole(roleWithPermissions);
      setShowPermissionsModal(true);
    } catch (error) {
      console.error('Ошибка загрузки роли:', error);
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1" gutterBottom>
          Управление ролями
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => setShowCreateModal(true)}
          sx={{ ml: 2 }}
        >
          Создать роль
        </Button>
      </Box>

      <Paper sx={{ width: '100%', overflow: 'hidden' }}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Описание</TableCell>
                <TableCell>Пользователи</TableCell>
                <TableCell>Дата создания</TableCell>
                <TableCell align="center">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {roles.map((role, index) => (
                <TableRow key={role.id || `role-${index}`} hover>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Security sx={{ mr: 1, color: 'primary.main' }} />
                      <Typography variant="body2" fontWeight="medium">
                        {role.name || '-'}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" color="text.secondary">
                      {role.description || '-'}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={roleUserCounts[role.id] || 0}
                      size="small"
                      color="primary"
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {role.created_at ? new Date(role.created_at).toLocaleDateString('ru-RU') : '-'}
                    </Typography>
                  </TableCell>
                  <TableCell align="center">
                    <Tooltip title="Редактировать">
                      <IconButton
                        size="small"
                        onClick={async () => {
                          try {
                            const roleWithPermissions = await rolesApi.getRole(role.id);
                            setSelectedRole(roleWithPermissions);
                            setShowEditModal(true);
                          } catch (error) {
                            console.error('Ошибка загрузки роли:', error);
                          }
                        }}
                        color="primary"
                      >
                        <Edit />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Права">
                      <IconButton
                        size="small"
                        onClick={() => handleViewPermissions(role)}
                        color="secondary"
                      >
                        <Security />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Удалить">
                      <IconButton
                        size="small"
                        onClick={() => handleDeleteRole(role.id)}
                        color="error"
                      >
                        <Delete />
                      </IconButton>
                    </Tooltip>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {showCreateModal && (
        <CreateRoleModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateRole}
          permissions={permissions}
        />
      )}

      {showEditModal && selectedRole && (
        <EditRoleModal
          role={selectedRole}
          onClose={() => {
            setShowEditModal(false);
            setSelectedRole(null);
          }}
          onSubmit={handleUpdateRole}
          permissions={permissions}
        />
      )}

      {showPermissionsModal && selectedRole && (
        <RolePermissionsModal
          role={selectedRole}
          onClose={() => {
            setShowPermissionsModal(false);
            setSelectedRole(null);
          }}
          permissions={permissions}
        />
      )}
    </Container>
  );
};

// Компонент создания роли
const CreateRoleModal: React.FC<{
  onClose: () => void;
  onSubmit: (data: any) => void;
  permissions: Permission[];
}> = ({ onClose, onSubmit, permissions }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    permission_ids: [] as string[],
    permissionSearch: ''
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}

    if (!formData.name.trim()) {
      errors.name = 'Название роли обязательно'
    }

    const filteredPermissionIds = formData.permission_ids.filter(id => id !== null && id !== undefined && id !== '');
    if (filteredPermissionIds.length === 0) {
      errors.permissions = 'Выберите хотя бы одно право'
    }

    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return

    try {
      setSubmitting(true)
      // Фильтруем null/undefined значения из permission_ids
      const filteredPermissionIds = formData.permission_ids.filter(id => id !== null && id !== undefined && id !== '');
      
      const submitData = {
        ...formData,
        permission_ids: filteredPermissionIds
      };
      
      console.log('Отправляемые данные роли:', submitData);
      await onSubmit(submitData);
    } finally {
      setSubmitting(false)
    }
  };

  const groupedPermissions = permissions.reduce((acc, perm) => {
    const module = perm.module || 'Общие';
    if (!acc[module]) {
      acc[module] = [];
    }
    acc[module].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <Dialog open onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Создать роль</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название роли"
            fullWidth
            variant="outlined"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={!!formErrors.name}
            helperText={formErrors.name}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Описание"
            fullWidth
            variant="outlined"
            multiline
            rows={3}
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            sx={{ mb: 2 }}
          />
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Права
            </Typography>
            {formErrors.permissions && (
              <Typography variant="caption" color="error">
                {formErrors.permissions}
              </Typography>
            )}
            <Box sx={{ maxHeight: 400, overflow: 'auto' }}>
              {Object.entries(groupedPermissions).map(([module, perms]) => (
                <Accordion key={module} defaultExpanded>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Box display="flex" alignItems="center" width="100%">
                      <Typography variant="subtitle1" sx={{ mr: 1 }}>
                        {getModuleIcon(module)}
                      </Typography>
                      <Typography variant="subtitle1" sx={{ flexGrow: 1 }}>
                        {module}
                      </Typography>
                      <FormControlLabel
                        control={
                          <Checkbox
                            checked={perms.every(perm => formData.permission_ids.includes(perm.id))}
                            onChange={(e) => {
                              if (e.target.checked) {
                                // Выбрать все права в модуле
                                const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                                const newPermissionIds = [...new Set([...formData.permission_ids, ...modulePermissionIds])];
                                setFormData({
                                  ...formData,
                                  permission_ids: newPermissionIds
                                });
                              } else {
                                // Снять все права в модуле
                                const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                                const newPermissionIds = formData.permission_ids.filter(id => !modulePermissionIds.includes(id));
                                setFormData({
                                  ...formData,
                                  permission_ids: newPermissionIds
                                });
                              }
                            }}
                            size="small"
                          />
                        }
                        label="Выбрать всё"
                        sx={{ m: 0 }}
                        onClick={(e) => e.stopPropagation()}
                      />
                    </Box>
                  </AccordionSummary>
                  <AccordionDetails>
                    <FormGroup>
                      {perms.map((perm, index) => {
                        const pid = perm.id;
                        const pcode = perm.code;
                        const pdesc = perm.description;
                        return (
                          <FormControlLabel
                            key={pid || `perm-${index}`}
                            control={
                              <Checkbox
                                checked={pid ? formData.permission_ids.includes(pid) : false}
                                onChange={(e) => {
                                  if (pid) {
                                    if (e.target.checked) {
                                      setFormData({
                                        ...formData,
                                        permission_ids: [...new Set([...formData.permission_ids, pid])]
                                      });
                                    } else {
                                      setFormData({
                                        ...formData,
                                        permission_ids: formData.permission_ids.filter(id => id !== pid)
                                      });
                                    }
                                  }
                                }}
                                size="small"
                              />
                            }
                            label={
                              <Box>
                                <Typography variant="body2">
                                  {pcode}
                                </Typography>
                                {pdesc && (
                                  <Typography variant="caption" color="text.secondary">
                                    {pdesc}
                                  </Typography>
                                )}
                              </Box>
                            }
                            sx={{ mb: 0.5 }}
                          />
                        );
                      })}
                    </FormGroup>
                  </AccordionDetails>
                </Accordion>
              ))}
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Отмена</Button>
          <Button
            type="submit"
            variant="contained"
            disabled={submitting}
          >
            {submitting ? <CircularProgress size={20} /> : 'Создать'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

// Компонент редактирования роли
const EditRoleModal: React.FC<{
  role: RoleWithPermissions;
  onClose: () => void;
  onSubmit: (id: string, data: any) => void;
  permissions: Permission[];
}> = ({ role, onClose, onSubmit, permissions }) => {

  const buildFormStateFromRole = () => {
    const rawName = role.name || (role as any).Name || '';
    const rawDescription = role.description || (role as any).Description || '';
    const permissionCodesRaw = role.permissions || (role as any).Permissions || [];
    const permissionCodes = Array.isArray(permissionCodesRaw) ? permissionCodesRaw : [];

    // Преобразуем коды прав в ID прав
    const permissionIds = permissionCodes
      .map(code => permissions.find(p => p.code === code)?.id)
      .filter((id): id is string => Boolean(id));

    console.log('buildFormStateFromRole - permissionCodes:', permissionCodes);
    console.log('buildFormStateFromRole - permissionIds:', permissionIds);
    console.log('buildFormStateFromRole - available permissions count:', permissions.length);

    return {
      name: rawName,
      description: rawDescription,
      permission_ids: permissionIds
    };
  };

  const [formData, setFormData] = useState(buildFormStateFromRole);
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    console.log('EditRoleModal useEffect triggered');
    console.log('role:', role);
    console.log('permissions count:', permissions.length);
    
    const nextFormData = buildFormStateFromRole();
    console.log('nextFormData:', nextFormData);

    // Всегда обновляем данные формы при изменении роли
    console.log('Updating form data with new data');
    setFormData(nextFormData);
  }, [role, permissions]);

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}

    if (!formData.name.trim()) {
      errors.name = 'Название роли обязательно'
    }

    const filteredPermissionIds = formData.permission_ids.filter(id => id !== null && id !== undefined && id !== '');
    if (filteredPermissionIds.length === 0) {
      errors.permissions = 'Выберите хотя бы одно право'
    }

    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return

    try {
      setSubmitting(true)
      const roleId = role.id || (role as any).ID;
      await onSubmit(roleId, formData);
    } finally {
      setSubmitting(false)
    }
  };

  const groupedPermissions = permissions.reduce((acc, perm) => {
    const module = perm.module || (perm as any).Module || 'Общие';
    if (!acc[module]) {
      acc[module] = [];
    }
    acc[module].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <Dialog open onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Редактировать роль</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название роли"
            fullWidth
            variant="outlined"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={!!formErrors.name}
            helperText={formErrors.name}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Описание"
            fullWidth
            variant="outlined"
            multiline
            rows={3}
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            sx={{ mb: 2 }}
          />
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Права
            </Typography>
            {formErrors.permissions && (
              <Typography variant="caption" color="error">
                {formErrors.permissions}
              </Typography>
            )}
            <Box sx={{ maxHeight: 400, overflow: 'auto' }}>
              {Object.entries(groupedPermissions).map(([module, perms]) => (
                <Accordion key={module} defaultExpanded>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Box display="flex" alignItems="center" width="100%">
                      <Typography variant="subtitle1" sx={{ mr: 1 }}>
                        {getModuleIcon(module)}
                      </Typography>
                      <Typography variant="subtitle1" sx={{ flexGrow: 1 }}>
                        {module}
                      </Typography>
                      <FormControlLabel
                        control={
                          <Checkbox
                            checked={perms.every(perm => formData.permission_ids.includes(perm.id))}
                            onChange={(e) => {
                              if (e.target.checked) {
                                // Выбрать все права в модуле
                                const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                                const newPermissionIds = [...new Set([...formData.permission_ids, ...modulePermissionIds])];
                                setFormData({
                                  ...formData,
                                  permission_ids: newPermissionIds
                                });
                              } else {
                                // Снять все права в модуле
                                const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                                const newPermissionIds = formData.permission_ids.filter(id => !modulePermissionIds.includes(id));
                                setFormData({
                                  ...formData,
                                  permission_ids: newPermissionIds
                                });
                              }
                            }}
                            size="small"
                          />
                        }
                        label="Выбрать всё"
                        sx={{ m: 0 }}
                        onClick={(e) => e.stopPropagation()}
                      />
                    </Box>
                  </AccordionSummary>
                  <AccordionDetails>
                    <FormGroup>
                      {perms.map((perm, index) => {
                        const pid = perm.id;
                        const pcode = perm.code;
                        const pdesc = perm.description;
                        return (
                          <FormControlLabel
                            key={pid || `perm-${index}`}
                            control={
                              <Checkbox
                                checked={pid ? formData.permission_ids.includes(pid) : false}
                                onChange={(e) => {
                                  if (pid) {
                                    if (e.target.checked) {
                                      setFormData({
                                        ...formData,
                                        permission_ids: [...new Set([...formData.permission_ids, pid])]
                                      });
                                    } else {
                                      setFormData({
                                        ...formData,
                                        permission_ids: formData.permission_ids.filter(id => id !== pid)
                                      });
                                    }
                                  }
                                }}
                                size="small"
                              />
                            }
                            label={
                              <Box>
                                <Typography variant="body2">
                                  {pcode}
                                </Typography>
                                {pdesc && (
                                  <Typography variant="caption" color="text.secondary">
                                    {pdesc}
                                  </Typography>
                                )}
                              </Box>
                            }
                            sx={{ mb: 0.5 }}
                          />
                        );
                      })}
                    </FormGroup>
                  </AccordionDetails>
                </Accordion>
              ))}
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Отмена</Button>
          <Button
            type="submit"
            variant="contained"
            disabled={submitting}
          >
            {submitting ? <CircularProgress size={20} /> : 'Сохранить'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

// Компонент просмотра прав роли
const RolePermissionsModal: React.FC<{
  role: RoleWithPermissions;
  onClose: () => void;
  permissions: Permission[];
}> = ({ role, onClose, permissions }) => {
  const groupedPermissions = permissions.reduce((acc, perm) => {
    const module = perm.module || 'Общие';
    if (!acc[module]) {
      acc[module] = [];
    }
    acc[module].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <Dialog open onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Права роли: {role.name}</DialogTitle>
      <DialogContent>
        <Box sx={{ maxHeight: 400, overflow: 'auto' }}>
          {Object.entries(groupedPermissions).map(([module, perms]) => (
            <Accordion key={module} defaultExpanded>
              <AccordionSummary expandIcon={<ExpandMore />}>
                <Box display="flex" alignItems="center">
                  <Typography variant="subtitle1" sx={{ mr: 1 }}>
                    {getModuleIcon(module)}
                  </Typography>
                  <Typography variant="subtitle1">
                    {module}
                  </Typography>
                </Box>
              </AccordionSummary>
              <AccordionDetails>
                <Box>
                  {perms.map((perm, index) => {
                    const hasPermission = role.permissions && Array.isArray(role.permissions) && role.permissions.includes(perm.code);
                    return (
                      <Box
                        key={perm.id || `perm-${index}`}
                        sx={{
                          display: 'flex',
                          alignItems: 'center',
                          p: 1,
                          mb: 0.5,
                          borderRadius: 1,
                          backgroundColor: hasPermission ? 'success.light' : 'grey.100',
                          color: hasPermission ? 'success.dark' : 'text.secondary'
                        }}
                      >
                        <Typography variant="body2" fontWeight="medium" sx={{ flexGrow: 1 }}>
                          {perm.code}
                        </Typography>
                        {perm.description && (
                          <Typography variant="caption" sx={{ ml: 1 }}>
                            - {perm.description}
                          </Typography>
                        )}
                        {hasPermission && (
                          <Typography variant="body2" color="success.main" sx={{ ml: 'auto' }}>
                            ✓
                          </Typography>
                        )}
                      </Box>
                    );
                  })}
                </Box>
              </AccordionDetails>
            </Accordion>
          ))}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Закрыть
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default RolesManagementPage;
