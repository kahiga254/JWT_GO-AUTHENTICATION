import { Box, Typography, useTheme } from "@mui/material";
import { NavLink } from "react-router-dom";
import { tokens } from "../../theme";

const Sidebar = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  return (
    <Box
      sx={{
        width: "250px",
        backgroundColor: colors.primary[400],
        padding: "20px",
        height: "100vh",
        position: "fixed",
        top: 0,
        left: 0,
      }}
    >
      <Typography variant="h4" color={colors.grey[100]} mb="20px">
        Admin Panel
      </Typography>
      <Box display="flex" flexDirection="column" gap="10px">
        <NavLink
          to="/"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Dashboard
        </NavLink>
        <NavLink
          to="/team"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Team
        </NavLink>
        <NavLink
          to="/contacts"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Contacts
        </NavLink>
        <NavLink
          to="/invoices"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Invoices
        </NavLink>
        <NavLink
          to="/form"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Form
        </NavLink>
        <NavLink
          to="/bar"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Bar
        </NavLink>
        <NavLink
          to="/pie"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Pie
        </NavLink>
        <NavLink
          to="/line"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Line
        </NavLink>
        <NavLink
          to="/faq"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          FAQ
        </NavLink>
        <NavLink
          to="/geography"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Geography
        </NavLink>
        <NavLink
          to="/calendar"
          style={({ isActive }) => ({
            color: isActive ? colors.greenAccent[500] : colors.grey[100],
            textDecoration: "none",
            padding: "5px 10px",
          })}
        >
          Calendar
        </NavLink>
      </Box>
    </Box>
  );
};

export default Sidebar;