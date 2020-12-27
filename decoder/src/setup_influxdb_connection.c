#include "aivdm.h"

static	void		set_client_socket_params(int socket_fd)
{
	t_timeval		timeout;
	int				on;

	timeout.tv_sec  = 0;
	timeout.tv_usec = 1;
	on = 1;
	setsockopt(socket_fd,  SOL_SOCKET, SO_REUSEADDR, (char *)&on, sizeof(on));
	setsockopt(socket_fd, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout,
															sizeof(timeout));
	return ;
}

static t_tls_session	*init_tls_session(t_tls_connection *tls_connection)
{
	t_tls_session		*tls_session;

	tls_session = (t_tls_session *)ft_memalloc(sizeof(*tls_session));
	tls_session->connection = tls_connection;
	tls_session->connection_status = e_waiting_msg2;
	return(tls_session);
}

t_tls_session			*setup_influxdb_connection(char *host_name,
															char *port_number)
{
	t_tls_connection	*tls_connection;
	SSL_CTX				*ctx;
	int					socket_fd;
	t_tls_session		*tls_session;
	char				read_buf[BUF_MAX_SIZE];
	int					chars;

	jk_init_openssl();
	tls_connection = NULL;
	tls_session = NULL;
	while (!tls_connection)
	{
		ctx = jk_start_tls_client(PEM_CERT_FILE, PEM_PRIVTE_KEY_FILE,
																	&socket_fd);
		tls_connection = jk_setup_tls_connection(host_name, port_number,
																socket_fd, ctx);
		set_client_socket_params(socket_fd);
		tls_session = init_tls_session(tls_connection);
		tls_session->connection_status = e_send_msg0;
	}
	while((chars = SSL_read(tls_connection->ssl_bio, read_buf, BUF_MAX_SIZE)) > 0)
			ft_printf("%s\n", read_buf);
	return(tls_session);
}
